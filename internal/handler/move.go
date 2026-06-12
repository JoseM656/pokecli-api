package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/JoseM656/pokecli-api/internal/cache"
	"github.com/JoseM656/pokecli-api/internal/model"
	"github.com/JoseM656/pokecli-api/internal/pokeapi"
	"github.com/go-chi/chi/v5"
)

// MoveHandler handles HTTP requests for Move data.
type MoveHandler struct {
	cache  cache.MoveRepository
	client *pokeapi.Client
}

// NewMoveHandler creates a new MoveHandler.
func NewMoveHandler(cache cache.MoveRepository, client *pokeapi.Client) *MoveHandler {
	return &MoveHandler{
		cache:  cache,
		client: client,
	}
}

// Get handles GET /move/{name}
func (h *MoveHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}

	// 1. Check cache
	move, err := h.cache.GetMoveByName(r.Context(), name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to check cache")
		return
	}

	// 2. Cache miss — fetch from PokeAPI
	if move == nil {
		raw, err := h.client.FetchMove(r.Context(), name)
		if err != nil {
			if errors.Is(err, model.ErrMoveNotFound) {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "move not found",
					"hint":  "try the English name, e.g. /move/flamethrower",
				})
				return
			}

			respondError(w, http.StatusInternalServerError, "failed to fetch move")
			return
		}

		move, err = h.client.MapMove(r.Context(), raw)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to map move")
			return
		}

		if err := h.cache.SaveMove(r.Context(), move); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to save to cache")
			return
		}

		move.Cached = false
	} else {
		move.Cached = true
	}

	move.Lang = lang
	respondJSON(w, http.StatusOK, move)
}
