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

func NewMatchMongoHandler(matches *repository.MatchRepo, teams *repository.TeamRepo, events chan<- models.EventLog) *MatchMongoHandler {
	return &MatchMongoHandler{matches: matches, teams: teams, events: events}
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

// Create match: matchKey auto-generated
func (h *MatchMongoHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DateTime string `json:"dateTime"` // RFC3339
		HomeCode string `json:"homeCode"`
		AwayCode string `json:"awayCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	req.HomeCode = strings.ToUpper(strings.TrimSpace(req.HomeCode))
	req.AwayCode = strings.ToUpper(strings.TrimSpace(req.AwayCode))
	req.DateTime = strings.TrimSpace(req.DateTime)

	if req.HomeCode == "" || req.AwayCode == "" || req.DateTime == "" {
		writeJSON(w, 400, map[string]string{"error": "dateTime, homeCode, awayCode required"})
		return
	}
	if req.HomeCode == req.AwayCode {
		writeJSON(w, 400, map[string]string{"error": "homeCode and awayCode must differ"})
		return
	}

	dt, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "dateTime must be RFC3339"})
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

	//  unique matchKey
	matchKey := req.HomeCode + "-" + req.AwayCode + "-" + time.Now().Format("20060102-150405")

	m := models.Match{
		MatchKey:  matchKey,
		DateTime:  dt,
		HomeCode:  req.HomeCode,
		AwayCode:  req.AwayCode,
		HomeGoals: 0,
		AwayGoals: 0,
		Status:    models.Scheduled,
		Events:    []models.MatchEvent{},
	}

	created, err := h.matches.Create(ctx, m)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "create error"})
		return
	}

	writeJSON(w, 201, created)
}

// PATCH /matches/{key}/events
func (h *MatchMongoHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.PathValue("key"))
	if key == "" {
		writeJSON(w, 400, map[string]string{"error": "missing match key"})
		return
	}

	var req models.MatchEvent
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	req.TeamCode = strings.ToUpper(strings.TrimSpace(req.TeamCode))
	req.Player = strings.TrimSpace(req.Player)
	req.Detail = strings.TrimSpace(req.Detail)
	req.CardColor = strings.ToLower(strings.TrimSpace(req.CardColor))
	req.PlayerOut = strings.TrimSpace(req.PlayerOut)
	req.PlayerIn = strings.TrimSpace(req.PlayerIn)

	if req.Type == "" || req.TeamCode == "" || req.Minute <= 0 || req.Minute > 130 {
		writeJSON(w, 400, map[string]string{"error": "type, teamCode, minute required (minute 1..130)"})
		return
	}

	allowed := map[string]bool{"goal": true, "card": true, "injury": true, "var": true, "sub": true}
	if !allowed[req.Type] {
		writeJSON(w, 400, map[string]string{"error": "type must be goal|card|injury|var|sub"})
		return
	}

	// validation by type
	if req.Type == "goal" {
		if req.Player == "" {
			writeJSON(w, 400, map[string]string{"error": "goal requires player"})
			return
		}
	}
	if req.Type == "card" {
		if req.Player == "" || (req.CardColor != "yellow" && req.CardColor != "red") {
			writeJSON(w, 400, map[string]string{"error": "card requires player and cardColor yellow|red"})
			return
		}
	}
	if req.Type == "sub" {
		if req.PlayerOut == "" || req.PlayerIn == "" {
			writeJSON(w, 400, map[string]string{"error": "sub requires playerOut and playerIn"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	// Ensure team exists
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
	if m.Status == models.Finished {
		writeJSON(w, 409, map[string]string{"error": "match already finished"})
		return
	}
	// Only allow events for teams playing in this match
	if req.TeamCode != m.HomeCode && req.TeamCode != m.AwayCode {
		writeJSON(w, 400, map[string]string{"error": "teamCode is not playing in this match"})
		return
	}

	//  repo allows status scheduled OR live (your $in filter)
	if err := h.matches.AddEvent(ctx, key, m, req); err != nil {
		writeJSON(w, 500, map[string]string{"error": "update error"})
		return
	}

	writeJSON(w, 200, map[string]string{"status": "ok"})
}

// PATCH /matches/{key}/status
func (h *MatchMongoHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.PathValue("key"))
	if key == "" {
		writeJSON(w, 400, map[string]string{"error": "missing match key"})
		return
	}

	var req struct {
		Status string `json:"status"` // scheduled | live | finished
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	s := models.MatchStatus(strings.ToLower(strings.TrimSpace(req.Status)))
	if s != models.Scheduled && s != models.Live && s != models.Finished {
		writeJSON(w, 400, map[string]string{"error": "status must be scheduled|live|finished"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	m, found, err := h.matches.FindByKey(ctx, key)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	if !found {
		writeJSON(w, 404, map[string]string{"error": "match not found"})
		return
	}

	// rule: finished cannot go back
	if m.Status == models.Finished && s != models.Finished {
		writeJSON(w, 409, map[string]string{"error": "finished match cannot be changed"})
		return
	}

	// if status finished -> use Finalize too (keeps rules consistent)
	if s == models.Finished {
		if err := h.matches.Finalize(ctx, key); err != nil {
			writeJSON(w, 500, map[string]string{"error": "finalize error"})
			return
		}
		writeJSON(w, 200, map[string]string{"status": "finished"})
		return
	}

	if err := h.matches.SetStatus(ctx, key, s); err != nil {
		writeJSON(w, 500, map[string]string{"error": "update error"})
		return
	}

	writeJSON(w, 200, map[string]string{"status": string(s)})
}

// POST /matches/{key}/finalize
func (h *MatchMongoHandler) Finalize(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.PathValue("key"))
	if key == "" {
		writeJSON(w, 400, map[string]string{"error": "missing match key"})
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

	if err := h.matches.Finalize(ctx, key); err != nil {
		writeJSON(w, 500, map[string]string{"error": "finalize error"})
		return
	}

	writeJSON(w, 200, map[string]string{"status": "finished"})
}
