package httpjson

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hferr/pack-api/internal/app"
)

func (h *Handler) ListPackSizes(w http.ResponseWriter, r *http.Request) {
	packs, err := h.packService.ListPacks(r.Context())
	if err != nil {
		http.Error(w, "Could not list packs", http.StatusInternalServerError)
		return
	}

	res := buildPacksSizeListResponse(packs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

type CreatePackRequest struct {
	Size int `json:"size"`
}

func (h *Handler) CreatePack(w http.ResponseWriter, r *http.Request) {
	var req CreatePackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Could not parse request body", http.StatusBadRequest)
		return
	}

	pack, err := h.packService.CreatePack(r.Context(), req.Size)
	if err != nil {
		var valErr *app.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Could not create pack", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pack)
}

type RebuildPacksRequest struct {
	Sizes []int `json:"sizes"`
}

func (h *Handler) RebuildPacks(w http.ResponseWriter, r *http.Request) {
	var req RebuildPacksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Could not parse request body", http.StatusBadRequest)
		return
	}

	packs, err := h.packService.RebuildPacks(r.Context(), req.Sizes)
	if err != nil {
		var valErr *app.ValidationError
		if errors.As(err, &valErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Could not rebuild packs", http.StatusInternalServerError)
		return
	}

	res := buildPacksSizeListResponse(packs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

type CalculatePacksForItemsRequest struct {
	Items int `json:"items"`
}

func (h *Handler) CalculateMinPackOrder(w http.ResponseWriter, r *http.Request) {
	var req CalculatePacksForItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Could not parse request body", http.StatusBadRequest)
		return
	}

	res, err := h.packService.CalculateMinPackOrder(r.Context(), req.Items)
	if err != nil {
		http.Error(w, "Error occurred while calculating order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// buildPacksSizeListResponse builds the reponse to be only an array of ints for simplicity
func buildPacksSizeListResponse(packs app.Packs) []int {
	res := make([]int, len(packs))
	for i, p := range packs {
		res[i] = p.Size
	}

	return res
}
