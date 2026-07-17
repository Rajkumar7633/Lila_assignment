package service

import (
	"context"

	"github.com/lila/player-viz-backend/internal/domain"
)

type MatchService struct {
	repo domain.Repository
	reg  *TransformerRegistry
}

func NewMatchService(repo domain.Repository, reg *TransformerRegistry) *MatchService {
	return &MatchService{repo: repo, reg: reg}
}

func (s *MatchService) ListDates(ctx context.Context) ([]string, error) {
	return s.repo.ListDates(ctx)
}

func (s *MatchService) ListMaps() []domain.MapInfo {
	return s.reg.ListMaps()
}

func (s *MatchService) ListMatches(ctx context.Context, filter domain.MatchFilter) ([]domain.MatchSummary, error) {
	return s.repo.ListMatches(ctx, filter)
}

func (s *MatchService) GetMatchJourney(ctx context.Context, matchID string) (*domain.MatchJourney, error) {
	journey, err := s.repo.GetMatchJourney(ctx, matchID)
	if err != nil {
		return nil, err
	}
	s.enrichJourneyMeta(journey)
	return s.transformJourney(journey)
}

func (s *MatchService) enrichJourneyMeta(journey *domain.MatchJourney) {
	for _, p := range journey.Paths {
		if p.IsBot {
			journey.BotCount++
		} else {
			journey.HumanCount++
		}
	}
	for _, ev := range journey.Events {
		switch ev.Category {
		case "kill":
			journey.EventCounts.Kills++
		case "death":
			journey.EventCounts.Deaths++
		case "loot":
			journey.EventCounts.Loot++
		case "storm":
			journey.EventCounts.Storm++
		}
	}
}

func (s *MatchService) transformJourney(journey *domain.MatchJourney) (*domain.MatchJourney, error) {
	mapID := journey.MapID

	for i := range journey.Paths {
		for j, pt := range journey.Paths[i].Points {
			pixel, err := s.reg.Transform(mapID, pt.X, pt.Y)
			if err != nil {
				return nil, err
			}
			journey.Paths[i].Points[j] = pixel
		}
	}

	for i, ev := range journey.Events {
		pixel, err := s.reg.Transform(mapID, ev.Position.X, ev.Position.Y)
		if err != nil {
			return nil, err
		}
		journey.Events[i].Position = pixel
	}

	return journey, nil
}
