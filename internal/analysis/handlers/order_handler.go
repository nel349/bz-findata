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

func (h *OrderHandler) StoreReceivedOrdersInSupabase(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseBodyParams(r)
	
	err := h.service.StoreReceivedOrdersInSupabase(r.Context(), hours, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "success")
}

func (h *OrderHandler) StoreMatchOrdersInSupabase(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseBodyParams(r)
	
	err := h.service.StoreMatchOrdersInSupabase(r.Context(), hours, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, "success")
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

func parseBodyParams(r *http.Request) (hours, limit int) {
	// Parse JSON body from request body
	var body struct {
		Hours int `json:"hours"`
		Limit int `json:"limit"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	return body.Hours, body.Limit
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}