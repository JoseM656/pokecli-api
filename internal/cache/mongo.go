package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pokeservicies/pokeproxy/internal/model"
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
	normalized := strings.ToLower(name)
	filter := bson.M{
		"$or": []bson.M{
			{"name.en": normalized},
			{"name.es": normalized},
			{"name.ja": normalized},
			{"name.fr": normalized},
			{"name.de": normalized},
			{"name.ko": normalized},
			{"name.it": normalized},
		},
	}

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

// MOVE --- ZONE

// MongoMoveRepository implements MoveRepository using MongoDB.
type MongoMoveRepository struct {
	collection *mongo.Collection
}

// NewMongoMoveRepository creates a new MongoMoveRepository.
func NewMongoMoveRepository(client *mongo.Client, dbName string) *MongoMoveRepository {
	collection := client.Database(dbName).Collection("moves")
	return &MongoMoveRepository{collection: collection}
}

// GetMoveByName looks up a Move by name in the cache.
// Returns nil, nil if the Move is not cached yet.
func (r *MongoMoveRepository) GetMoveByName(ctx context.Context, name string) (*model.Move, error) {
	normalized := strings.ToLower(name)
	filter := bson.M{
		"$or": []bson.M{
			{"name.en": normalized},
			{"name.es": normalized},
			{"name.ja": normalized},
			{"name.fr": normalized},
			{"name.de": normalized},
			{"name.ko": normalized},
			{"name.it": normalized},
		},
	}

	var move model.Move
	err := r.collection.FindOne(ctx, filter).Decode(&move)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get move from cache: %w", err)
	}

	return &move, nil
}

// SaveMove stores a Move in the cache.
func (r *MongoMoveRepository) SaveMove(ctx context.Context, move *model.Move) error {
	filter := bson.M{"name.en": move.Name.EN}
	update := bson.M{"$set": move}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save move to cache: %w", err)
	}

	return nil
}
