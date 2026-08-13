package main

import (
	"log"
	"os"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	log.Printf("Config loaded successfully")
	log.Println("Configs", cfg)
}
