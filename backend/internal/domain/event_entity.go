package domain

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type GameEvent struct {
	UserID    string `json:"userId"`
	IsBot     bool   `json:"isBot"`
	EventType string `json:"eventType"`
	Category  string `json:"category"`
	TS        int64  `json:"ts"`
	Position  Point  `json:"position"`
}

type PlayerPath struct {
	UserID string  `json:"userId"`
	IsBot  bool    `json:"isBot"`
	Points []Point `json:"points"`
	Times  []int64 `json:"times"`
}

type HeatmapType string

const (
	HeatmapKills   HeatmapType = "kills"
	HeatmapDeaths  HeatmapType = "deaths"
	HeatmapTraffic HeatmapType = "traffic"
)

type HeatmapFilter struct {
	MapID     string
	DateLabel string
	MatchID   string
	Type      HeatmapType
}

type HeatmapPoint struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Intensity float64 `json:"intensity"`
}

type RawHeatmapRow struct {
	X         float64
	Z         float64
	EventType string
	MapID     string
}
