package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // Required: Postgres driver import
)

// Client wraps the standard SQL DB connection.
// This is the struct passed around the application when database access is needed.
type Client struct {
	DB *sql.DB
}

// NewClient initializes and returns a new Postgres database client.
// It uses the databaseURL provided in the configuration.
func NewClient(databaseURL string) *Client {
	// Open the database connection using the postgres driver
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Set connection pool limits for a highly scalable app
	db.SetMaxOpenConns(25)              // Max number of open connections
	db.SetMaxIdleConns(5)               // Max number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Max time a connection can be reused

	// Ping the database to ensure connection is valid
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database (Supabase): %v", err)
	}

	log.Println("Successfully connected to the PostgreSQL database (Supabase).")

	return &Client{
		DB: db,
	}
}

// Close closes the database connection pool.
func (c *Client) Close() {
	if c.DB != nil {
		c.DB.Close()
		log.Println("Database connection pool closed.")
	}
}