package handlers

import (
	"context"
	"net/http"
	"sort"
	"time"

	"final-by-me/internal/repository"
)

type TableRow struct {
	TeamCode string `json:"teamCode"`
	TeamName string `json:"teamName"`
	P        int    `json:"played"`
	W        int    `json:"wins"`
	D        int    `json:"draws"`
	L        int    `json:"losses"`
	GF       int    `json:"goalsFor"`
	GA       int    `json:"goalsAgainst"`
	GD       int    `json:"goalDiff"`
	Pts      int    `json:"points"`
}

type TableHandler struct {
	teams   *repository.TeamRepo
	matches *repository.MatchRepo
}

func NewTableHandler(teams *repository.TeamRepo, matches *repository.MatchRepo) *TableHandler {
	return &TableHandler{teams: teams, matches: matches}
}

func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	teams, err := h.teams.List(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}

	rows := make(map[string]*TableRow, len(teams))
	for _, t := range teams {
		rows[t.Code] = &TableRow{TeamCode: t.Code, TeamName: t.Name}
	}

	finished, err := h.matches.ListFinished(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}

	for _, m := range finished {
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

		// 3/1/0
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
			home.Pts += 1
			away.Pts += 1
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

	writeJSON(w, 200, map[string]any{"table": out, "count": len(out)})
}
