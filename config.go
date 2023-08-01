package main

import "github.com/kelseyhightower/envconfig"

type (
	// Config provides the system configuration.
	Config struct {
		Server Server
	}

	Server struct {
		Port string `default:"8080"`
	}
)

// Environ returns the settings from the environment.
func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}
