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
	mux.HandleFunc("POST /packs/rebuild", h.RebuildPacks)
	mux.HandleFunc("POST /packs/calculate-order", h.CalculateMinPackOrder)

	handler := corsMiddleware(mux)
	return handler
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
