package main

import (
	"log"

	"helm-viewer/config"
	"helm-viewer/router"
)

func main() {
	cfg := config.NewConfig()

	r := router.SetupRouter()

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
