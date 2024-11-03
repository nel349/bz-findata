package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nel349/bz-findata/internal/analysis"
)

type OrderHandler struct {
	service *analysis.Service
}

func NewOrderHandler(service *analysis.Service) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetLargestReceivedOrders(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseQueryParams(r, 2, 10)
	
	orders, err := h.service.GetLargestReceivedOrdersInLastNHours(r.Context(), hours, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, orders)
}

func (h *OrderHandler) GetLargestOpenOrders(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseQueryParams(r, 24, 100)
	
	orders, err := h.service.GetLargestOpenOrdersInLastNHours(r.Context(), hours, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, orders)
}

func (h *OrderHandler) GetLargestMatchOrders(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseQueryParams(r, 24, 100)
	
	orders, err := h.service.GetLargestMatchOrdersInLastNHours(r.Context(), hours, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, orders)
}

// Helper functions
func parseQueryParams(r *http.Request, defaultHours, defaultLimit int) (hours, limit int) {
	hours = defaultHours
	limit = defaultLimit

	if h := r.URL.Query().Get("hours"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil {
			hours = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	return hours, limit
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}