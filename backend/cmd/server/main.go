package main

import (
	"log"

	"github.com/lila/player-viz-backend/internal/config"
	"github.com/lila/player-viz-backend/internal/repository"
	"github.com/lila/player-viz-backend/internal/service"
	httptransport "github.com/lila/player-viz-backend/internal/transport/http"
)

func main() {
	cfg := config.Load()

	repo, err := repository.NewDuckDBRepository(cfg.ParquetGlob())
	if err != nil {
		log.Fatalf("repository init: %v", err)
	}
	defer repo.Close()

	reg := service.NewTransformerRegistry()
	matchSvc := service.NewMatchService(repo, reg)
	heatmapSvc := service.NewHeatmapService(repo, reg)

	router := httptransport.NewRouter(cfg, matchSvc, heatmapSvc)
	addr := ":" + cfg.Port
	log.Printf("server listening on %s (data: %s)", addr, cfg.DataPath)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
