import { useFilters } from '../../context/FilterContext';
import type { HeatmapType } from '../../types';

const HEATMAP_OPTIONS: { id: HeatmapType; label: string }[] = [
  { id: 'traffic', label: 'Traffic' },
  { id: 'kills', label: 'Kills' },
  { id: 'deaths', label: 'Deaths' },
];

export function MapSelector({ maps }: { maps: { id: string; label: string }[] }) {
  const { mapId, setMapId, setMatchId } = useFilters();

  return (
    <label className="field">
      <span>Map</span>
      <select
        value={mapId}
        onChange={(e) => {
          setMapId(e.target.value);
          setMatchId('');
        }}
      >
        <option value="">Choose map…</option>
        {maps.map((m) => (
          <option key={m.id} value={m.id}>
            {m.label}
          </option>
        ))}
      </select>
    </label>
  );
}

export function DatePicker({ dates }: { dates: string[] }) {
  const { dateLabel, setDateLabel, setMatchId } = useFilters();

  return (
    <label className="field">
      <span>Date</span>
      <select
        value={dateLabel}
        onChange={(e) => {
          setDateLabel(e.target.value);
          setMatchId('');
        }}
      >
        <option value="">All dates</option>
        {dates.map((d) => (
          <option key={d} value={d}>
            {d.replace('_', ' ')}
          </option>
        ))}
      </select>
    </label>
  );
}

const formatDuration = (ms: number): string => {
  const totalSeconds = Math.round(ms / 1000);
  const m = Math.floor(totalSeconds / 60);
  const s = totalSeconds % 60;
  return `${m}m ${s}s`;
};

export function MatchSelector({
  matches,
  loading,
}: {
  matches: { matchId: string; humanCount: number; botCount: number; durationMs: number }[];
  loading: boolean;
}) {
  const { matchId, setMatchId } = useFilters();

  return (
    <label className="field field-wide">
      <span>Match</span>
      <select
        value={matchId}
        disabled={loading || matches.length === 0}
        onChange={(e) => setMatchId(e.target.value)}
      >
        <option value="">
          {loading ? 'Loading matches…' : matches.length ? 'Choose match…' : 'No matches'}
        </option>
        {matches.map((m) => (
          <option key={m.matchId} value={m.matchId}>
            {m.matchId.slice(0, 8)}… · {m.humanCount}H / {m.botCount}B · {formatDuration(m.durationMs)}
          </option>
        ))}
      </select>
    </label>
  );
}

export function HeatmapToggle() {
  const { heatmapType, setHeatmapType } = useFilters();

  return (
    <div className="toggle-group">
      <span className="toggle-label">Heatmap</span>
      <button
        type="button"
        className={heatmapType == null ? 'chip active' : 'chip'}
        onClick={() => setHeatmapType(null)}
      >
        Off
      </button>
      {HEATMAP_OPTIONS.map((opt) => (
        <button
          key={opt.id}
          type="button"
          className={heatmapType === opt.id ? 'chip active' : 'chip'}
          onClick={() => setHeatmapType(opt.id)}
        >
          {opt.label}
        </button>
      ))}
    </div>
  );
}

export function PlayerToggle() {
  const { showHumans, showBots, setShowHumans, setShowBots } = useFilters();

  return (
    <div className="toggle-group">
      <span className="toggle-label">Show</span>
      <button
        type="button"
        className={showHumans ? 'chip active human' : 'chip'}
        onClick={() => setShowHumans(!showHumans)}
      >
        Humans
      </button>
      <button
        type="button"
        className={showBots ? 'chip active bot' : 'chip'}
        onClick={() => setShowBots(!showBots)}
      >
        Bots
      </button>
    </div>
  );
}
