package handlers

import (
	"net/http"

	"github.com/nel349/bz-findata/internal/analysis/dex"
)

type DexHandler struct {
	service *dex.Service
}

func NewDexHandler(service *dex.Service) *DexHandler {
	return &DexHandler{
		service: service,
	}
}

// Get the largest swaps in last N hours by Value
func (h *DexHandler) GetLargestSwaps(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseQueryParams(r, 24, 100)
	swaps, err := h.service.GetLargestSwapsInLastNHours(r.Context(), hours, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, swaps)
}
