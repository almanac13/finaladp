package handlers

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"final-by-me/internal/repository"
)

type TableRow struct {
	TeamCode   string `json:"teamCode"`
	TeamName   string `json:"teamName"`
	League     string `json:"league"`
	IsFavorite bool   `json:"isFavorite"`

	P   int `json:"played"`
	W   int `json:"wins"`
	D   int `json:"draws"`
	L   int `json:"losses"`
	GF  int `json:"goalsFor"`
	GA  int `json:"goalsAgainst"`
	GD  int `json:"goalDiff"`
	Pts int `json:"points"`
}

type TableHandler struct {
	teams   *repository.TeamRepo
	matches *repository.MatchRepo
}

func NewTableHandler(teams *repository.TeamRepo, matches *repository.MatchRepo) *TableHandler {
	return &TableHandler{teams: teams, matches: matches}
}

// GET /table?league=EPL&favorite=ARS
// If league is provided -> table for that league only.
// Favorite is optional -> marks that team row as isFavorite=true.
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	league := strings.TrimSpace(r.URL.Query().Get("league"))
	fav := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("favorite")))

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()

	// Load teams (by league if specified)
	var teamsList []struct {
		Code   string
		Name   string
		League string
	}

	if league == "" {
		all, err := h.teams.List(ctx)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": "db error"})
			return
		}
		for _, t := range all {
			teamsList = append(teamsList, struct {
				Code, Name, League string
			}{t.Code, t.Name, t.League})
		}
	} else {
		list, err := h.teams.ListByLeague(ctx, league)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": "db error"})
			return
		}
		for _, t := range list {
			teamsList = append(teamsList, struct {
				Code, Name, League string
			}{t.Code, t.Name, t.League})
		}
	}

	// Build rows
	rows := make(map[string]*TableRow, len(teamsList))
	leagueSet := make(map[string]bool, len(teamsList))
	for _, t := range teamsList {
		leagueSet[t.Code] = true
		rows[t.Code] = &TableRow{
			TeamCode:   t.Code,
			TeamName:   t.Name,
			League:     t.League,
			IsFavorite: (fav != "" && t.Code == fav),
		}
	}

	// Finished matches
	finished, err := h.matches.ListFinished(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}

	for _, m := range finished {
		// If league specified: only count matches where BOTH teams are in this league
		if league != "" {
			if !leagueSet[m.HomeCode] || !leagueSet[m.AwayCode] {
				continue
			}
		}

		home := rows[m.HomeCode]
		away := rows[m.AwayCode]
		if home == nil || away == nil {
			continue
		}

		home.P++
		away.P++

		home.GF += m.HomeGoals
		home.GA += m.AwayGoals
		away.GF += m.AwayGoals
		away.GA += m.HomeGoals

		if m.HomeGoals > m.AwayGoals {
			home.W++
			away.L++
			home.Pts += 3
		} else if m.HomeGoals < m.AwayGoals {
			away.W++
			home.L++
			away.Pts += 3
		} else {
			home.D++
			away.D++
			home.Pts++
			away.Pts++
		}
	}

	out := make([]TableRow, 0, len(rows))
	for _, r := range rows {
		r.GD = r.GF - r.GA
		out = append(out, *r)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Pts != out[j].Pts {
			return out[i].Pts > out[j].Pts
		}
		if out[i].GD != out[j].GD {
			return out[i].GD > out[j].GD
		}
		if out[i].GF != out[j].GF {
			return out[i].GF > out[j].GF
		}
		return out[i].TeamName < out[j].TeamName
	})

	writeJSON(w, 200, map[string]any{
		"league": league,
		"count":  len(out),
		"table":  out,
	})
}
