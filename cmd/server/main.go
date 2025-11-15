package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mitcheltastic/ManproBackend/config"
	fbclient "github.com/mitcheltastic/ManproBackend/internal/infrastructure/firebase"
	"github.com/mitcheltastic/ManproBackend/internal/infrastructure/router"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Firebase Admin SDK
	log.Println("Initializing Firebase Admin SDK...")
	// The Firebase client is initialized with the service account JSON file path from config
	firebaseClient := fbclient.NewClient(cfg.FirebaseServiceKeyJSON)
	log.Println("Firebase Admin SDK initialized successfully.")

	// 3. Initialize HTTP Router
	log.Println("Initializing HTTP router (Gin)...")
	r := gin.Default()

	// 4. Set up the routes, passing the Firebase client to the router setup
	router.SetupRoutes(r, firebaseClient)

	// 5. Start the Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
