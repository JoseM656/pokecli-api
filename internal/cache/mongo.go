package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JoseM656/pokecli-api/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository implements Repository using MongoDB.
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a new MongoRepository.
func NewMongoRepository(client *mongo.Client, dbName string) *MongoRepository {
	collection := client.Database(dbName).Collection("pokemon")
	return &MongoRepository{collection: collection}
}

// Connect establishes a connection to MongoDB and verifies it with a ping.
func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client, nil
}

// GetByName looks up a Pokemon by name in the cache.
// Returns nil, nil if the Pokemon is not cached yet.
func (r *MongoRepository) GetByName(ctx context.Context, name string) (*model.Pokemon, error) {
	filter := bson.M{"name.en": name}

	var pokemon model.Pokemon
	err := r.collection.FindOne(ctx, filter).Decode(&pokemon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pokemon from cache: %w", err)
	}

	return &pokemon, nil
}

// Save stores a Pokemon in the cache.
func (r *MongoRepository) Save(ctx context.Context, pokemon *model.Pokemon) error {
	filter := bson.M{"_id": pokemon.ID}
	update := bson.M{"$set": pokemon}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save pokemon to cache: %w", err)
	}

	return nil
}
