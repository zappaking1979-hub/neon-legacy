package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/domain/player"
	"github.com/neonlegacy/server/internal/middleware"
)

type GymHandler struct {
	gymService *application.GymService
	playerRepo player.Repository
	tmpl       *template.Template
}

func NewGymHandler(gymService *application.GymService, playerRepo player.Repository, tmpl *template.Template) *GymHandler {
	return &GymHandler{gymService: gymService, playerRepo: playerRepo, tmpl: tmpl}
}

type gymPageData struct {
	Player    *player.Player
	Exercises []gymListEntry
}

type gymListEntry struct {
	ID          int
	Name        string
	Description string
	Stat        string
	EnergyCost  int
	MinLevel    int
	GainMin     int
	GainMax     int
	CanTrain    bool
	BlockedBy   string
}

func (h *GymHandler) Page(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)

	exercises, err := h.gymService.ListExercises(r.Context())
	if err != nil {
		http.Error(w, "failed to load exercises", http.StatusInternalServerError)
		return
	}

	entries := make([]gymListEntry, 0, len(exercises))
	for _, ex := range exercises {
		entry := gymListEntry{
			ID:          ex.ID,
			Name:        ex.Name,
			Description: ex.Description,
			Stat:        string(ex.Stat),
			EnergyCost:  ex.EnergyCost,
			MinLevel:    ex.MinLevel,
			GainMin:     ex.GainMin,
			GainMax:     ex.GainMax,
		}
		ok, reason := ex.CanTrain(p.Level, p.Energy)
		entry.CanTrain = ok
		if !ok {
			entry.BlockedBy = reason
		}
		entries = append(entries, entry)
	}

	h.tmpl.ExecuteTemplate(w, "pages/gym.html", gymPageData{
		Player:    p,
		Exercises: entries,
	})
}

func (h *GymHandler) Train(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	exerciseID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "invalid exercise id", http.StatusBadRequest)
		return
	}

	result, err := h.gymService.Train(r.Context(), p, exerciseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p, _ = h.playerRepo.GetByID(r.Context(), p.ID)

	h.tmpl.ExecuteTemplate(w, "partials/gym_result.html", map[string]interface{}{
		"Success":  result.Success,
		"StatGain": result.StatGain,
		"StatName": result.StatName,
		"Message":  result.Message,
		"ExpGain":  result.ExpGain,
		"Player":   p,
	})
}
