package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration settings.
type Config struct {
	// Path to the Firebase service account JSON file.
	// Comes from .env: FIREBASE_SERVICE_KEY_JSON=./scripts/firebase-key.json
	FirebaseServiceKeyJSON string `envconfig:"FIREBASE_SERVICE_KEY_JSON" required:"true"`

	// Server settings (e.g., port)
	Port string `envconfig:"PORT" default:"8080"`

	// Supabase/Postgres connection string
	// You can leave this empty for now or set a placeholder in .env.
	DatabaseURL string `envconfig:"DATABASE_URL" default:""`
}

// LoadConfig reads configuration from .env file and environment variables.
func LoadConfig() *Config {
	var cfg Config

	// Load .env if present (falls back to OS env if not found).
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or failed to load. Using system environment variables.")
	}

	// Populate cfg from environment variables.
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return &cfg
}
