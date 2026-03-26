package httpjson

import (
	"net/http"

	"github.com/hferr/pack-api/internal/app"
)

type Handler struct {
	packService app.PackService
}

func NewHandler(packService app.PackService) *Handler {
	return &Handler{
		packService: packService,
	}
}

func (h *Handler) NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /packs/sizes", h.ListPackSizes)
	mux.HandleFunc("POST /packs", h.CreatePack)
	mux.HandleFunc("POST /packs/calculate-order", h.CalculateMinPackOrder)

	return mux
}
