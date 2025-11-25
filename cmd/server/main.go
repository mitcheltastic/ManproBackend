package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/mitcheltastic/ManproBackend/config"
	// **FIXED ALIAS:** Renamed to 'dbimpl' to correctly access the package
	dbimpl "github.com/mitcheltastic/ManproBackend/internal/infrastructure/database" 
	fbclient "github.com/mitcheltastic/ManproBackend/internal/infrastructure/firebase"
	router "github.com/mitcheltastic/ManproBackend/internal/infrastructure/router"
)

// runMigrations executes the 'Up' migrations using the Goose library.
func runMigrations(db *sql.DB) {
	// IMPORTANT: Goose expects the migration directory path.
	// We use the absolute path to ensure it works regardless of where 'go run' is executed.
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b) 

	// Navigate up one level to the project root, then down to the migrations folder
	migrationsDir := filepath.Join(filepath.Dir(basepath), "..", "scripts", "migrations")
	
	log.Printf("Starting database migration from directory: %s", migrationsDir)
	
	// Configure Goose to use the standard Postgres dialect
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose failed to set dialect: %v", err)
	}

	// Run all pending 'Up' migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database migration completed successfully.")
}

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Firebase Admin SDK (Auth is already wired)
	log.Println("Initializing Firebase Admin SDK...")
	// *** CRITICAL FIX HERE ***: Use the path variable, not the old JSON variable
	firebaseClient := fbclient.NewClient(cfg.FirebaseServiceKeyPath)
	log.Println("Firebase Admin SDK initialized successfully.")

	// 3. Initialize Supabase (Postgres) Database Client
	log.Println("Initializing Supabase Database Client...")
	// **FIXED CALL:** Using the alias 'dbimpl'
	dbClient := dbimpl.NewClient(cfg.DatabaseURL)
	defer dbClient.Close()

	// 4. Run Migrations Conditionally
	if cfg.ShouldMigrate {
		// Only run if the environment flag is explicitly set to true
		runMigrations(dbClient.DB)
	} else {
		log.Println("Skipping automatic database migration (SHOULD_MIGRATE=false)")
	}
	
	// 5. Initialize HTTP Router
	log.Println("Initializing HTTP router (Gin)...")
	r := gin.Default()
	
	// We need to pass the database client to the router so handlers can access it
	router.SetupRoutes(r, firebaseClient, dbClient) 

	// 6. Start the Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}