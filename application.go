package main

import (
	"log/slog"
	"net/http"
	"time"
)

type application struct {
	settings           *settings
	logger             *slog.Logger
	handler            *http.ServeMux
	replayNonceService *InMemoryReplayNonceService
}

func newApplication(settings *settings, logger *slog.Logger) (*application, error) {
	app := &application{
		settings:           settings,
		logger:             logger,
		handler:            http.NewServeMux(),
		replayNonceService: NewInMemoryReplayNonceService(time.Hour, time.Hour),
	}

	app.handler.HandleFunc("GET /acme/directory", app.handleGetDirectory)
	app.handler.HandleFunc("GET /acme/newNonce", app.handleGetNonce)
	app.handler.HandleFunc("HEAD /acme/newNonce", app.handleHeadNonce)

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
