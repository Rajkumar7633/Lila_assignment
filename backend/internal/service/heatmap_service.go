package service

import (
	"context"
	"math"

	"github.com/lila/player-viz-backend/internal/domain"
)

type HeatmapService struct {
	repo domain.Repository
	reg  *TransformerRegistry
}

func NewHeatmapService(repo domain.Repository, reg *TransformerRegistry) *HeatmapService {
	return &HeatmapService{repo: repo, reg: reg}
}

const heatmapGridSize = 32

func (s *HeatmapService) GetHeatmap(ctx context.Context, filter domain.HeatmapFilter) ([]domain.HeatmapPoint, error) {
	rows, err := s.repo.GetHeatmapRaw(ctx, filter)
	if err != nil {
		return nil, err
	}

	type key struct{ gx, gy int }
	buckets := make(map[key]float64)
	maxCount := 0.0

	cell := float64(1024) / float64(heatmapGridSize)

	for _, row := range rows {
		pixel, err := s.reg.Transform(row.MapID, row.X, row.Z)
		if err != nil {
			return nil, err
		}
		gx := int(math.Floor(pixel.X / cell))
		gy := int(math.Floor(pixel.Y / cell))
		if gx < 0 || gy < 0 || gx >= heatmapGridSize || gy >= heatmapGridSize {
			continue
		}
		k := key{gx, gy}
		buckets[k]++
		if buckets[k] > maxCount {
			maxCount = buckets[k]
		}
	}

	if maxCount == 0 {
		return []domain.HeatmapPoint{}, nil
	}

	points := make([]domain.HeatmapPoint, 0, len(buckets))
	for k, count := range buckets {
		centerX := (float64(k.gx) + 0.5) * cell
		centerY := (float64(k.gy) + 0.5) * cell
		points = append(points, domain.HeatmapPoint{
			X:         centerX,
			Y:         centerY,
			Intensity: count / maxCount,
		})
	}
	return points, nil
}
