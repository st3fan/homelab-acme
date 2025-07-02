package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type GetDirectoryResponse struct {
	NewNonce   string `json:"newNonce"`
	NewAccount string `json:"newAccount"`
	NewOrder   string `json:"newOrder"`
	RevokeCert string `json:"revokeCert"`
	KeyChange  string `json:"keyChange"`
}

func (app *application) handleGetDirectory(w http.ResponseWriter, r *http.Request) {
	getDirectoryResponse := GetDirectoryResponse{
		NewNonce:   app.settings.ServerURL.ResolveReference(&url.URL{Path: "/acme/new-nonce"}).String(),
		NewAccount: app.settings.ServerURL.ResolveReference(&url.URL{Path: "/acme/new-account"}).String(),
		KeyChange:  app.settings.ServerURL.ResolveReference(&url.URL{Path: "/acme/key-change"}).String(),
		NewOrder:   app.settings.ServerURL.ResolveReference(&url.URL{Path: "/acme/new-order"}).String(),
		RevokeCert: app.settings.ServerURL.ResolveReference(&url.URL{Path: "/acme/revoke-cert"}).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(getDirectoryResponse); err != nil {
		app.logger.Error("Failed to encode JSON response", err)
	}
}

func (app *application) handleGetNonce(w http.ResponseWriter, _ *http.Request) {
	replayNonce, err := app.replayNonceService.New()
	if err != nil {
		app.logger.Error("Failed to generate nonce", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Replay-Nonce", replayNonce)
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handleHeadNonce(w http.ResponseWriter, _ *http.Request) {
	replayNonce, err := app.replayNonceService.New()
	if err != nil {
		app.logger.Error("Failed to generate nonce", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Replay-Nonce", replayNonce)
	w.WriteHeader(http.StatusOK)
}
