package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration settings.
type Config struct {
	// Firebase service account file path, used to initialize the Admin SDK.
	FirebaseServiceKeyPath string `envconfig:"FIREBASE_SERVICE_KEY_PATH" required:"true"`
	
	// Server settings (e.g., port)
	Port string `envconfig:"PORT" default:"8080"`

	// Supabase/Postgres connection string
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

// LoadConfig reads configuration from .env file and environment variables.
func LoadConfig() *Config {
	// Load .env file first (will be overridden by system environment variables)
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or failed to load. Using system environment variables.")
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return &cfg
}