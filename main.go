package main

import (
	"log/slog"
	"os"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	settings, err := newSettingsFromEnv()
	if err != nil {
		log.Error("Failed to load settings", "error", err)
		os.Exit(1)
	}

	app, err := newApplication(settings, log)
	if err != nil {
		log.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	if err := app.run(); err != nil {
		log.Error("Failed to run application", "error", err)
		os.Exit(1)
	}
}
