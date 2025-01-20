package handlers

import (
	"log"
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

// Store the largest swaps in last N hours by Value
func (h *DexHandler) StoreLargestSwaps(w http.ResponseWriter, r *http.Request) {
	hours, limit := parseBodyParams(r)
	// Log the parameters for debugging
	log.Printf("Storing largest swaps with hours: %d, limit: %d", hours, limit)
	err := h.service.StoreLargestSwapsInLastNHours(r.Context(), hours, limit)
	if err != nil {
		log.Println("error storing largest swaps", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Log success message
	log.Println("Successfully stored largest swaps")
	respondWithJSON(w, "success")
}
