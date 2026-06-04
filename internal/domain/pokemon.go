package domain

import "context"

// Pokemon represents the domain model of a Pokemon.
type Pokemon struct {
	ID        int      `json:"id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`
	Height    int      `json:"height" bson:"height"`
	Weight    int      `json:"weight" bson:"weight"`
	Types     []string `json:"types" bson:"types"`
	CachedAt  int64    `json:"cached_at,omitempty" bson:"cached_at"`
}

// PokemonRepository defines the database actions for Pokemon cache.
type PokemonRepository interface {
	GetByName(ctx context.Context, name string) (*Pokemon, error)
	Save(ctx context.Context, pokemon *Pokemon) error
}

// PokeAPIClient defines the external API communication actions.
type PokeAPIClient interface {
	FetchByName(ctx context.Context, name string) (*Pokemon, error)
}

// PokemonService defines the business logic operations.
type PokemonService interface {
	GetPokemon(ctx context.Context, name string) (*Pokemon, error)
}
