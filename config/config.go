package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration.
type Config struct {
	Port        string
	MongoURI    string
	MongoDBName string
	PokeAPIURL  string
}

// Load loads the configuration from environment variables or .env file.
func Load() *Config {
	// Attempt to load .env file; ignore if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// if fails use this standar configuration.
	return &Config{
		Port:        getEnv("PORT", "8080"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27107"),
		MongoDBName: getEnv("MONGO_DB_NAME", "pokecli"),
		PokeAPIURL:  getEnv("POKEAPI_URL", "https://pokeapi.co/api/v2"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
