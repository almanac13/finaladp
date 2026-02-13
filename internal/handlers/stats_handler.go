package handlers

import (
	"context"
	"net/http"
	"time"

	"final-by-me/internal/repository"
)

type StatsHandler struct {
	matches *repository.MatchRepo
}

func NewStatsHandler(matches *repository.MatchRepo) *StatsHandler {
	return &StatsHandler{matches: matches}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	all, err := h.matches.List(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}

	finished, err := h.matches.ListFinished(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}

	totalGoals := 0
	for _, m := range finished {
		totalGoals += m.HomeGoals + m.AwayGoals
	}

	avg := 0.0
	if len(finished) > 0 {
		avg = float64(totalGoals) / float64(len(finished))
	}

	writeJSON(w, 200, map[string]any{
		"totalMatches":     len(all),
		"finishedMatches":  len(finished),
		"totalGoals":       totalGoals,
		"avgGoalsPerMatch": avg,
	})
}
