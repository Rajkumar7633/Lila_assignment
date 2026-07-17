package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lila/player-viz-backend/internal/config"
	"github.com/lila/player-viz-backend/internal/middleware"
	"github.com/lila/player-viz-backend/internal/service"
	"github.com/lila/player-viz-backend/internal/transport/http/handlers"
)

func NewRouter(cfg config.Config, matchSvc *service.MatchService, heatmapSvc *service.HeatmapService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	matchHandler := handlers.NewMatchHandler(matchSvc)
	heatmapHandler := handlers.NewHeatmapHandler(heatmapSvc)

	r.GET("/health", handlers.Health)
	r.Static("/minimaps", cfg.MinimapPath)

	api := r.Group("/api")
	{
		api.GET("/maps", matchHandler.ListMaps)
		api.GET("/dates", matchHandler.ListDates)
		api.GET("/matches", matchHandler.ListMatches)
		api.GET("/matches/:matchId/journey", matchHandler.GetMatchJourney)
		api.GET("/heatmap", heatmapHandler.GetHeatmap)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return r
}
