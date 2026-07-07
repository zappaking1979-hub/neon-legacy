package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/domain/player"
	"github.com/neonlegacy/server/internal/middleware"
)

type InventoryHandler struct {
	inventoryService *application.InventoryService
	playerRepo       player.Repository
	tmpl             *template.Template
}

func NewInventoryHandler(inventoryService *application.InventoryService, playerRepo player.Repository, tmpl *template.Template) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService, playerRepo: playerRepo, tmpl: tmpl}
}

func (h *InventoryHandler) Page(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)

	items, err := h.inventoryService.ListInventory(r.Context(), p.ID)
	if err != nil {
		http.Error(w, "failed to load inventory", http.StatusInternalServerError)
		return
	}

	type invItemEntry struct {
		Item
		Quantity int
	}
	entries := make([]invItemEntry, 0, len(items))
	for _, it := range items {
		entries = append(entries, invItemEntry{
			Item:     it.Item,
			Quantity: it.Quantity,
		})
	}

	h.tmpl.ExecuteTemplate(w, "pages/inventory.html", map[string]interface{}{
		"Player": p,
		"Items":  entries,
	})
}

func (h *InventoryHandler) Use(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	msg, err := h.inventoryService.UseItem(r.Context(), p, itemID)
	message := ""
	success := false

	if err != nil {
		message = err.Error()
	} else {
		success = true
		message = msg
	}

	p, _ = h.playerRepo.GetByID(r.Context(), p.ID)

	h.tmpl.ExecuteTemplate(w, "partials/inventory_result.html", map[string]interface{}{
		"Success": success,
		"Message": message,
		"Player":  p,
	})
}