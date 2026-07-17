package domain

import "context"

// Repository abstracts parquet access (Dependency Inversion — services depend on this interface).
type Repository interface {
	ListDates(ctx context.Context) ([]string, error)
	ListMatches(ctx context.Context, filter MatchFilter) ([]MatchSummary, error)
	GetMatchJourney(ctx context.Context, matchID string) (*MatchJourney, error)
	GetHeatmapRaw(ctx context.Context, filter HeatmapFilter) ([]RawHeatmapRow, error)
}

// CoordinateTransformer maps world (x,z) to minimap pixel space (Strategy pattern).
type CoordinateTransformer interface {
	MapID() string
	WorldToPixel(x, z float64) (pixelX, pixelY float64)
	Config() MapInfo
}
