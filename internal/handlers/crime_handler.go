package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/domain/player"
	"github.com/neonlegacy/server/internal/middleware"
)

type CrimeHandler struct {
	crimeService *application.CrimeService
	playerRepo   player.Repository
	tmpl         *template.Template
}

func NewCrimeHandler(crimeService *application.CrimeService, playerRepo player.Repository, tmpl *template.Template) *CrimeHandler {
	return &CrimeHandler{crimeService: crimeService, playerRepo: playerRepo, tmpl: tmpl}
}

type crimesPageData struct {
	Player *player.Player
	Crimes []*CrimeListEntry
}

type CrimeListEntry struct {
	ID        int
	Name      string
	NerveCost int
	CanAfford bool
	CanCommit bool
	BlockedBy string
}

type crimeResultData struct {
	Player         *player.Player
	Success        bool
	Jailed         bool
	Message        string
	ExpGain        int
	CashGain       int
}

func (h *CrimeHandler) Page(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)
	if p == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	crimes, err := h.crimeService.ListCrimes(r.Context())
	if err != nil {
		http.Error(w, "failed to load crimes", http.StatusInternalServerError)
		return
	}

	var entries []*CrimeListEntry
	for _, c := range crimes {
		entry := &CrimeListEntry{
			ID:        c.ID,
			Name:      c.Name,
			NerveCost: c.NerveCost,
			CanAfford: p.Nerve >= c.NerveCost,
			CanCommit: true,
		}
		if p.Level < c.MinLevel {
			entry.CanCommit = false
			entry.BlockedBy = "Level " + strconv.Itoa(c.MinLevel)
		}
		entries = append(entries, entry)
	}

	h.tmpl.ExecuteTemplate(w, "crimes.html", crimesPageData{
		Player: p,
		Crimes: entries,
	})
}

func (h *CrimeHandler) DoCrime(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)
	if p == nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	crimeID, _ := strconv.Atoi(r.FormValue("id"))
	multiplier, _ := strconv.Atoi(r.FormValue("cm"))

	if crimeID < 1 {
		http.Error(w, "invalid crime", http.StatusBadRequest)
		return
	}

	switch {
	case multiplier < 1:
		multiplier = 1
	case multiplier > 10:
		multiplier = 10
	}

	result, err := h.crimeService.DoCrime(r.Context(), p, crimeID, multiplier)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	// Refresh player from DB for fresh stats
	p, _ = h.playerRepo.GetByID(r.Context(), p.ID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "crime_result.html", crimeResultData{
		Player:   p,
		Success:  result.Success,
		Jailed:   result.Jailed,
		Message:  result.Message,
		ExpGain:  result.ExpGain,
		CashGain: result.CashGain,
	})
}
