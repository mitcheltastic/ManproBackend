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

	// Flag to determine if database migrations should run on startup.
	ShouldMigrate bool `envconfig:"SHOULD_MIGRATE" default:"false"`

	// Secret key used to sign and verify local JWTs for custom authentication.
	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`

	// --- NEW EMAIL CONFIG ---
	SMTPHost string `envconfig:"SMTP_HOST" default:""`
	SMTPPort string `envconfig:"SMTP_PORT" default:""`
	SMTPUser string `envconfig:"SMTP_USER" default:""`
	SMTPPass string `envconfig:"SMTP_PASS" default:""`
	FromEmail string `envconfig:"FROM_EMAIL" default:""`
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
		log.Fatalf("Error loading configuration: required key missing value. %v", err)
	}

	return &cfg
}