package main

import (
	"log/slog"
	"net/http"
)

type application struct {
	settings *settings
	logger   *slog.Logger
	handler  *http.ServeMux
}

func newApplication(settings *settings, logger *slog.Logger) (*application, error) {
	app := &application{
		settings: settings,
		logger:   logger,
		handler:  http.NewServeMux(),
	}

	app.handler.HandleFunc("GET /acme/directory", app.handleGetDirectory)

	return app, nil
}

func (app *application) run() error {
	app.logger.Info("Starting application")

	httpServer := &http.Server{
		Addr:    app.settings.ServerAddress,
		Handler: app.handler,
	}

	return httpServer.ListenAndServe()
}
