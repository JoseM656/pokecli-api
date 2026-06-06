package cache

import (
	"context"

	"github.com/JoseM656/pokecli-api/internal/model"
)

// Repository defines the cache operations for Pokemon data.
type Repository interface {
	GetByName(ctx context.Context, name string) (*model.Pokemon, error)
	Save(ctx context.Context, pokemon *model.Pokemon) error
}
