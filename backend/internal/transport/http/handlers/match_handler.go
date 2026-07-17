package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lila/player-viz-backend/internal/domain"
	"github.com/lila/player-viz-backend/internal/service"
)

type MatchHandler struct {
	matches *service.MatchService
}

func NewMatchHandler(matches *service.MatchService) *MatchHandler {
	return &MatchHandler{matches: matches}
}

func (h *MatchHandler) ListMaps(c *gin.Context) {
	c.JSON(http.StatusOK, h.matches.ListMaps())
}

func (h *MatchHandler) ListDates(c *gin.Context) {
	dates, err := h.matches.ListDates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dates)
}

func (h *MatchHandler) ListMatches(c *gin.Context) {
	filter := domain.MatchFilter{
		MapID:     c.Query("map"),
		DateLabel: c.Query("date"),
	}
	matches, err := h.matches.ListMatches(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

func (h *MatchHandler) GetMatchJourney(c *gin.Context) {
	matchID := c.Param("matchId")
	journey, err := h.matches.GetMatchJourney(c.Request.Context(), matchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, journey)
}
