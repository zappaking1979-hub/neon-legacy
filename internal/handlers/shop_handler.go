package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/domain/item"
	"github.com/neonlegacy/server/internal/domain/player"
	"github.com/neonlegacy/server/internal/middleware"
)

type ShopHandler struct {
	shopService *application.ShopService
	playerRepo  player.Repository
	tmpl        *template.Template
}

func NewShopHandler(shopService *application.ShopService, playerRepo player.Repository, tmpl *template.Template) *ShopHandler {
	return &ShopHandler{shopService: shopService, playerRepo: playerRepo, tmpl: tmpl}
}

func (h *ShopHandler) Page(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)

	items, err := h.shopService.ListItems(r.Context())
	if err != nil {
		http.Error(w, "failed to load items", http.StatusInternalServerError)
		return
	}

	owned, _ := h.shopService.ListInventory(r.Context(), p.ID)
	ownedMap := make(map[int]int)
	for _, oi := range owned {
		ownedMap[oi.Item.ID] = oi.Quantity
	}

	type shopItemEntry struct {
		item.Item
		Owned int
	}
	entries := make([]shopItemEntry, 0, len(items))
	for _, it := range items {
		entries = append(entries, shopItemEntry{
			Item:  it,
			Owned: ownedMap[it.ID],
		})
	}

	h.tmpl.ExecuteTemplate(w, "pages/shop.html", map[string]interface{}{
		"Player": p,
		"Items":  entries,
	})
}

func (h *ShopHandler) Buy(w http.ResponseWriter, r *http.Request) {
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

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity <= 0 {
		quantity = 1
	}

	err = h.shopService.BuyItem(r.Context(), p, itemID, quantity)
	message := ""
	success := false

	if err != nil {
		message = err.Error()
	} else {
		success = true
		it, _ := h.shopService.GetItem(r.Context(), itemID)
		if it != nil {
			message = "Bought " + it.Name + " x" + strconv.Itoa(quantity) + " for $" + strconv.FormatInt(it.BuyPrice*int64(quantity), 10)
		} else {
			message = "Purchase successful"
		}
	}

	p, _ = h.playerRepo.GetByID(r.Context(), p.ID)

	h.tmpl.ExecuteTemplate(w, "partials/shop_result.html", map[string]interface{}{
		"Success": success,
		"Message": message,
		"Player":  p,
	})
}

func (h *ShopHandler) Sell(w http.ResponseWriter, r *http.Request) {
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

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity <= 0 {
		quantity = 1
	}

	err = h.shopService.SellItem(r.Context(), p, itemID, quantity)
	message := ""
	success := false

	if err != nil {
		message = err.Error()
	} else {
		success = true
		it, _ := h.shopService.GetItem(r.Context(), itemID)
		if it != nil {
			message = "Sold " + it.Name + " x" + strconv.Itoa(quantity) + " for $" + strconv.FormatInt(it.SellPrice*int64(quantity), 10)
		} else {
			message = "Sale successful"
		}
	}

	p, _ = h.playerRepo.GetByID(r.Context(), p.ID)

	h.tmpl.ExecuteTemplate(w, "partials/shop_result.html", map[string]interface{}{
		"Success": success,
		"Message": message,
		"Player":  p,
	})
}