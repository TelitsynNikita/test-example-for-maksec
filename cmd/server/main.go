package main

import (
	"log"
	"os"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/config"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	log.Printf("Config loaded successfully")
	log.Println("Configs", cfg)

	script := domain.NewScript(
		"192.168.1.10",
		"root",
		"template1",
		"/opt/script-monitor/scripts/123e4567-e89b-12d3-a456-426614174000.sh",
	)
	log.Printf("   Script created: ID=%s, Path=%s", script.ID, script.Path)

	eventTime, _ := time.Parse(time.RFC3339, "2026-08-13T10:00:00Z")
	event := domain.NewEvent(
		script.ID,
		"root",
		script.Path,
		"execute",
		eventTime,
	)
	log.Printf("   Event created: ID=%s, Action=%s", event.ID, event.Action)

	templates := domain.GetTemplates()
	log.Printf("   Available templates: %d", len(templates))
	for name, tmpl := range templates {
		log.Printf("     - %s: %s", name, tmpl.Description)
	}
}
