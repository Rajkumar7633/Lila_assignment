package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lila/player-viz-backend/internal/domain"
	"github.com/lila/player-viz-backend/internal/service"
)

type HeatmapHandler struct {
	heatmaps *service.HeatmapService
}

func NewHeatmapHandler(heatmaps *service.HeatmapService) *HeatmapHandler {
	return &HeatmapHandler{heatmaps: heatmaps}
}

func (h *HeatmapHandler) GetHeatmap(c *gin.Context) {
	filter := domain.HeatmapFilter{
		MapID:     c.Query("map"),
		DateLabel: c.Query("date"),
		MatchID:   c.Query("match"),
		Type:      domain.HeatmapType(c.DefaultQuery("type", "traffic")),
	}
	if filter.MapID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "map query param required"})
		return
	}

	points, err := h.heatmaps.GetHeatmap(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, points)
}
