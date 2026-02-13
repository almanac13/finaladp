package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"final-by-me/internal/repository"
	"final-by-me/internal/seed"
)

type TeamHandler struct {
	teams *repository.TeamRepo
}

func NewTeamHandler(teams *repository.TeamRepo) *TeamHandler {
	return &TeamHandler{teams: teams}
}

// GET /teams?league=EPL
func (h *TeamHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	league := strings.TrimSpace(r.URL.Query().Get("league"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var list any
	var err error

	if league == "" {
		list, err = h.teams.List(ctx)
	} else {
		list, err = h.teams.ListByLeague(ctx, league)
	}

	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	// list is []models.Team
	writeJSON(w, 200, map[string]any{"teams": list})
}

// GET /leagues
func (h *TeamHandler) ListLeagues(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"leagues": seed.Leagues()})
}
