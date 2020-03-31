package config

import "github.com/LaurenceGA/go-crev/internal/id"

type Configuration struct {
	CurrentStore string `yaml:"current_store,omitempty"`
	CurrentID    *id.ID `yaml:"current_id,omitempty"`
}
