package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"final-by-me/internal/models"
	"final-by-me/internal/repository"
)

type MatchMongoHandler struct {
	matches *repository.MatchRepo
	teams   *repository.TeamRepo
	events  chan<- models.EventLog
}

func NewMatchMongoHandler(
	matches *repository.MatchRepo,
	teams *repository.TeamRepo,
	events chan<- models.EventLog,
) *MatchMongoHandler {
	return &MatchMongoHandler{
		matches: matches,
		teams:   teams,
		events:  events,
	}
}

func isAdmin(r *http.Request) bool {
	return strings.ToLower(strings.TrimSpace(r.Header.Get("X-Role"))) == "admin"
}

func (h *MatchMongoHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	list, err := h.matches.List(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"matches": list, "count": len(list)})
}

func (h *MatchMongoHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		writeJSON(w, 403, map[string]string{"error": "admin only (X-Role: admin)"})
		return
	}

	var req struct {
		MatchKey string `json:"matchKey"`
		DateTime string `json:"dateTime"` // ISO string
		HomeCode string `json:"homeCode"`
		AwayCode string `json:"awayCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	req.MatchKey = strings.TrimSpace(req.MatchKey)
	req.HomeCode = strings.ToUpper(strings.TrimSpace(req.HomeCode))
	req.AwayCode = strings.ToUpper(strings.TrimSpace(req.AwayCode))

	if req.MatchKey == "" || req.HomeCode == "" || req.AwayCode == "" || req.DateTime == "" {
		writeJSON(w, 400, map[string]string{"error": "matchKey, dateTime, homeCode, awayCode required"})
		return
	}

	dt, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "dateTime must be RFC3339, e.g. 2026-02-10T18:30:00Z"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	homeOK, err := h.teams.Exists(ctx, req.HomeCode)
	if err != nil || !homeOK {
		writeJSON(w, 400, map[string]string{"error": "homeCode not found"})
		return
	}
	awayOK, err := h.teams.Exists(ctx, req.AwayCode)
	if err != nil || !awayOK {
		writeJSON(w, 400, map[string]string{"error": "awayCode not found"})
		return
	}
	if req.HomeCode == req.AwayCode {
		writeJSON(w, 400, map[string]string{"error": "homeCode and awayCode must differ"})
		return
	}

	m := models.Match{
		MatchKey:  req.MatchKey,
		DateTime:  dt,
		HomeCode:  req.HomeCode,
		AwayCode:  req.AwayCode,
		HomeGoals: 0,
		AwayGoals: 0,
		Status:    models.Scheduled,
		Goals:     []models.Goal{},
		Cards:     []models.Card{},
	}

	created, err := h.matches.Create(ctx, m)
	if err != nil {
		// duplicate key => matchKey exists
		writeJSON(w, 409, map[string]string{"error": "matchKey already exists"})
		return
	}
	writeJSON(w, 201, created)
}

func (h *MatchMongoHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		writeJSON(w, 403, map[string]string{"error": "admin only (X-Role: admin)"})
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/matches/")
	key = strings.TrimSuffix(key, "/events")
	key = strings.TrimSpace(key)

	if key == "" {
		writeJSON(w, 400, map[string]string{"error": "missing match key"})
		return
	}

	var req struct {
		Type     string `json:"type"` // goal | card
		TeamCode string `json:"teamCode"`
		Player   string `json:"player"`
		Minute   int    `json:"minute"`
		Color    string `json:"color"` // for card
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	req.TeamCode = strings.ToUpper(strings.TrimSpace(req.TeamCode))
	req.Player = strings.TrimSpace(req.Player)
	req.Color = strings.ToLower(strings.TrimSpace(req.Color))

	if req.Type != "goal" && req.Type != "card" {
		writeJSON(w, 400, map[string]string{"error": "type must be 'goal' or 'card'"})
		return
	}
	if req.TeamCode == "" || req.Player == "" || req.Minute <= 0 || req.Minute > 130 {
		writeJSON(w, 400, map[string]string{"error": "teamCode, player, minute required (minute 1..130)"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	ok, err := h.teams.Exists(ctx, req.TeamCode)
	if err != nil || !ok {
		writeJSON(w, 400, map[string]string{"error": "teamCode not found"})
		return
	}

	// Ensure match exists
	m, found, err := h.matches.FindByKey(ctx, key)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	if !found {
		writeJSON(w, 404, map[string]string{"error": "match not found"})
		return
	}
	// Only allow events for teams participating
	if req.TeamCode != m.HomeCode && req.TeamCode != m.AwayCode {
		writeJSON(w, 400, map[string]string{"error": "teamCode is not playing in this match"})
		return
	}
	if m.Status != models.Scheduled {
		writeJSON(w, 409, map[string]string{"error": "match is not scheduled"})
		return
	}

	if req.Type == "goal" {
		err = h.matches.AddGoal(ctx, key, models.Goal{
			TeamCode: req.TeamCode,
			Player:   req.Player,
			Minute:   req.Minute,
		})
	} else {
		if req.Color != "yellow" && req.Color != "red" {
			writeJSON(w, 400, map[string]string{"error": "card color must be 'yellow' or 'red'"})
			return
		}
		err = h.matches.AddCard(ctx, key, models.Card{
			TeamCode: req.TeamCode,
			Player:   req.Player,
			Color:    req.Color,
			Minute:   req.Minute,
		})
	}

	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "update error"})
		return
	}

	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *MatchMongoHandler) Finalize(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		writeJSON(w, 403, map[string]string{"error": "admin only (X-Role: admin)"})
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/matches/")
	key = strings.TrimSuffix(key, "/finalize")
	key = strings.TrimSpace(key)

	var req struct {
		HomeGoals int `json:"homeGoals"`
		AwayGoals int `json:"awayGoals"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.HomeGoals < 0 || req.AwayGoals < 0 {
		writeJSON(w, 400, map[string]string{"error": "goals must be >= 0"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	_, found, err := h.matches.FindByKey(ctx, key)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	if !found {
		writeJSON(w, 404, map[string]string{"error": "match not found"})
		return
	}

	if err := h.matches.Finalize(ctx, key, req.HomeGoals, req.AwayGoals); err != nil {
		writeJSON(w, 500, map[string]string{"error": "finalize error"})
		return
	}

	writeJSON(w, 200, map[string]string{"status": "finished"})
}
