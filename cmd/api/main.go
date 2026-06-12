package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JoseM656/pokecli-api/config"
	"github.com/JoseM656/pokecli-api/internal/cache"
	"github.com/JoseM656/pokecli-api/internal/handler"
	"github.com/JoseM656/pokecli-api/internal/pokeapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Config
	cfg := config.Load()

	// MongoDB
	mongoClient, err := cache.Connect(context.Background(), cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Dependencies
	pokeRepo := cache.NewMongoRepository(mongoClient, cfg.MongoDBName)
	moveRepo := cache.NewMongoMoveRepository(mongoClient, cfg.MongoDBName)
	pokeClient := pokeapi.NewClient(cfg.PokeAPIURL)

	pokemonHandler := handler.NewPokemonHandler(pokeRepo, pokeClient)
	moveHandler := handler.NewMoveHandler(moveRepo, pokeClient)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status": "ok"}`)
	})

	r.Get("/pokemon/{name}", pokemonHandler.Get)
	r.Get("/move/{name}", moveHandler.Get)

	// Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server running on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
