package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/JoseM656/pokecli-api/internal/cache"
	"github.com/JoseM656/pokecli-api/internal/pokeapi"
	"github.com/go-chi/chi/v5"

	"github.com/JoseM656/pokecli-api/internal/model"
)

// PokemonHandler handles HTTP requests for Pokemon data.
type PokemonHandler struct {
	cache  cache.Repository
	client *pokeapi.Client
}

// NewPokemonHandler creates a new PokemonHandler.
func NewPokemonHandler(cache cache.Repository, client *pokeapi.Client) *PokemonHandler {
	return &PokemonHandler{
		cache:  cache,
		client: client,
	}
}

// Get handles GET /pokemon/{name}
func (h *PokemonHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}

	// 1. Check cache
	pokemon, err := h.cache.GetByName(r.Context(), name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to check cache")
		return
	}

	// 2. Cache miss — fetch from PokeAPI
	if pokemon == nil {
		raw, err := h.client.FetchPokemon(r.Context(), name)
		if err != nil {
			if errors.Is(err, model.ErrPokemonNotFound) {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "pokemon not found",
					"hint":  "try the English name, e.g. /pokemon/pikachu",
				})
				return
			}
			respondError(w, http.StatusInternalServerError, "failed to fetch pokemon")
			return
		}

		pokemon, err = h.client.Map(r.Context(), raw)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to map pokemon")
			return
		}

		if err := h.cache.Save(r.Context(), pokemon); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to save to cache")
			return
		}

		pokemon.Cached = false
	} else {
		pokemon.Cached = true
	}

	pokemon.Lang = lang
	respondJSON(w, http.StatusOK, pokemon)
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
