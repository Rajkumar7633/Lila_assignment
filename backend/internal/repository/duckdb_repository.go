package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/lila/player-viz-backend/internal/domain"
)

type DuckDBRepository struct {
	db          *sql.DB
	parquetGlob string
}

func NewDuckDBRepository(parquetGlob string) (*DuckDBRepository, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("open duckdb: %w", err)
	}
	return &DuckDBRepository{db: db, parquetGlob: parquetGlob}, nil
}

func (r *DuckDBRepository) Close() error {
	return r.db.Close()
}

func (r *DuckDBRepository) ListDates(ctx context.Context) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT regexp_extract(filename, '(February_[0-9]+)', 1) AS date_label
		FROM read_parquet('%s', filename=true)
		WHERE date_label IS NOT NULL AND date_label != ''
		ORDER BY date_label
	`, r.parquetGlob)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list dates: %w", err)
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}
	return dates, rows.Err()
}

func (r *DuckDBRepository) ListMatches(ctx context.Context, filter domain.MatchFilter) ([]domain.MatchSummary, error) {
	query := fmt.Sprintf(`
		SELECT
			match_id,
			map_id,
			MIN(regexp_extract(filename, '(February_[0-9]+)', 1)) AS date_label,
			COUNT(DISTINCT user_id) AS player_count,
			MIN(epoch_ms(ts)) AS start_ts,
			MAX(epoch_ms(ts)) AS end_ts
		FROM read_parquet('%s', filename=true)
		WHERE ($1 = '' OR map_id = $1)
		  AND ($2 = '' OR regexp_extract(filename, '(February_[0-9]+)', 1) = $2)
		GROUP BY match_id, map_id
		ORDER BY start_ts DESC
	`, r.parquetGlob)

	rows, err := r.db.QueryContext(ctx, query, filter.MapID, filter.DateLabel)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()

	var matches []domain.MatchSummary
	for rows.Next() {
		var m domain.MatchSummary
		if err := rows.Scan(&m.MatchID, &m.MapID, &m.DateLabel, &m.PlayerCount, &m.StartTS, &m.EndTS); err != nil {
			return nil, err
		}
		m.DurationMS = (m.EndTS - m.StartTS) * 1000
		matches = append(matches, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return matches, nil
	}

	if err := r.enrichHumanBotCounts(ctx, matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *DuckDBRepository) enrichHumanBotCounts(ctx context.Context, matches []domain.MatchSummary) error {
	ids := make([]string, len(matches))
	for i, m := range matches {
		ids[i] = m.MatchID
	}
	inClause := buildInClause(ids)

	query := fmt.Sprintf(`
		SELECT match_id, user_id
		FROM read_parquet('%s', filename=true)
		WHERE match_id IN (%s)
		GROUP BY match_id, user_id
	`, r.parquetGlob, inClause)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	counts := make(map[string][2]int)
	for rows.Next() {
		var matchID, userID string
		if err := rows.Scan(&matchID, &userID); err != nil {
			return err
		}
		c := counts[matchID]
		if domain.IsBot(userID) {
			c[1]++
		} else {
			c[0]++
		}
		counts[matchID] = c
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for i := range matches {
		c := counts[matches[i].MatchID]
		matches[i].HumanCount = c[0]
		matches[i].BotCount = c[1]
	}
	return nil
}

func (r *DuckDBRepository) GetMatchJourney(ctx context.Context, matchID string) (*domain.MatchJourney, error) {
	query := fmt.Sprintf(`
		SELECT
			user_id,
			match_id,
			map_id,
			x,
			z,
			epoch_ms(ts) AS ts_ms,
			CAST(event AS VARCHAR) AS event_type,
			regexp_extract(filename, '(February_[0-9]+)', 1) AS date_label
		FROM read_parquet('%s', filename=true)
		WHERE match_id = $1
		ORDER BY ts_ms ASC
	`, r.parquetGlob)

	rows, err := r.db.QueryContext(ctx, query, matchID)
	if err != nil {
		return nil, fmt.Errorf("get match journey: %w", err)
	}
	defer rows.Close()

	journey := &domain.MatchJourney{MatchID: matchID}
	pathIndex := make(map[string]int)

	for rows.Next() {
		var userID, mid, mapID, eventType, dateLabel string
		var x, z float64
		var ts int64
		if err := rows.Scan(&userID, &mid, &mapID, &x, &z, &ts, &eventType, &dateLabel); err != nil {
			return nil, err
		}

		if journey.MapID == "" {
			journey.MapID = mapID
			journey.DateLabel = dateLabel
		}
		if journey.StartTS == 0 || ts < journey.StartTS {
			journey.StartTS = ts
		}
		if ts > journey.EndTS {
			journey.EndTS = ts
		}

		isBot := domain.IsBot(userID)

		if domain.IsMovementEvent(eventType) {
			idx, ok := pathIndex[userID]
			if !ok {
				journey.Paths = append(journey.Paths, domain.PlayerPath{
					UserID: userID,
					IsBot:  isBot,
				})
				idx = len(journey.Paths) - 1
				pathIndex[userID] = idx
			}
			journey.Paths[idx].Points = append(journey.Paths[idx].Points, domain.Point{X: x, Y: z})
			journey.Paths[idx].Times = append(journey.Paths[idx].Times, ts)
			continue
		}

		journey.Events = append(journey.Events, domain.GameEvent{
			UserID:    userID,
			IsBot:     isBot,
			EventType: eventType,
			Category:  domain.EventCategory(eventType),
			TS:        ts,
			Position:  domain.Point{X: x, Y: z},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if journey.MapID == "" {
		return nil, fmt.Errorf("match not found: %s", matchID)
	}
	return journey, nil
}

func (r *DuckDBRepository) GetHeatmapRaw(ctx context.Context, filter domain.HeatmapFilter) ([]domain.RawHeatmapRow, error) {
	eventFilter := heatmapEventFilter(filter.Type)
	query := fmt.Sprintf(`
		SELECT x, z, CAST(event AS VARCHAR) AS event_type, map_id
		FROM read_parquet('%s', filename=true)
		WHERE map_id = $1
		  AND ($2 = '' OR regexp_extract(filename, '(February_[0-9]+)', 1) = $2)
		  AND ($3 = '' OR match_id = $3)
		  AND CAST(event AS VARCHAR) IN (%s)
	`, r.parquetGlob, eventFilter)

	rows, err := r.db.QueryContext(ctx, query, filter.MapID, filter.DateLabel, filter.MatchID)
	if err != nil {
		return nil, fmt.Errorf("heatmap raw: %w", err)
	}
	defer rows.Close()

	var out []domain.RawHeatmapRow
	for rows.Next() {
		var row domain.RawHeatmapRow
		if err := rows.Scan(&row.X, &row.Z, &row.EventType, &row.MapID); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func heatmapEventFilter(t domain.HeatmapType) string {
	switch t {
	case domain.HeatmapKills:
		return "'Kill','BotKill'"
	case domain.HeatmapDeaths:
		return "'Killed','BotKilled','KilledByStorm'"
	default:
		return "'Position','BotPosition'"
	}
}

func buildInClause(values []string) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	}
	return strings.Join(parts, ",")
}
