package config

type Configuration struct {
	CurrentStore string `yaml:"current_store,omitempty"`
	CurrentID    string `yaml:"current_id,omitempty"`
}
