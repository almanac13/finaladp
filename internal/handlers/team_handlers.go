package handlers

import (
	"context"
	"net/http"
	"time"

	"final-by-me/internal/repository"
)

type TeamHandler struct {
	teams *repository.TeamRepo
}

func NewTeamHandler(teams *repository.TeamRepo) *TeamHandler {
	return &TeamHandler{teams: teams}
}

func (h *TeamHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	list, err := h.teams.List(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"teams": list, "count": len(list)})
}
