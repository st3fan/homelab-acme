package main

import (
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type settings struct {
	Domain        string  `default:"example.com"`
	ServerAddress string  `default:"127.0.0.1:8080"`
	ServerURL     url.URL `default:"https://acme.example.com:8080"`
}

func newSettingsFromEnv() (*settings, error) {
	var settings settings
	if err := envconfig.Process("HOMELAB_ACME", &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}
