package domain

type MapInfo struct {
	ID          string  `json:"id"`
	Label       string  `json:"label"`
	Scale       float64 `json:"scale"`
	OriginX     float64 `json:"originX"`
	OriginZ     float64 `json:"originZ"`
	ImageURL    string  `json:"imageUrl"`
	ImageWidth  int     `json:"imageWidth"`
	ImageHeight int     `json:"imageHeight"`
}

type MatchSummary struct {
	MatchID     string `json:"matchId"`
	MapID       string `json:"mapId"`
	DateLabel   string `json:"dateLabel"`
	PlayerCount int    `json:"playerCount"`
	HumanCount  int    `json:"humanCount"`
	BotCount    int    `json:"botCount"`
	StartTS     int64  `json:"startTs"`
	EndTS       int64  `json:"endTs"`
	DurationMS  int64  `json:"durationMs"`
}

type MatchFilter struct {
	MapID     string
	DateLabel string
}

type MatchJourney struct {
	MatchID     string       `json:"matchId"`
	MapID       string       `json:"mapId"`
	DateLabel   string       `json:"dateLabel"`
	StartTS     int64        `json:"startTs"`
	EndTS       int64        `json:"endTs"`
	HumanCount  int          `json:"humanCount"`
	BotCount    int          `json:"botCount"`
	EventCounts EventCounts  `json:"eventCounts"`
	Paths       []PlayerPath `json:"paths"`
	Events      []GameEvent  `json:"events"`
}

type EventCounts struct {
	Kills  int `json:"kills"`
	Deaths int `json:"deaths"`
	Loot   int `json:"loot"`
	Storm  int `json:"storm"`
}
