package trust

import (
	"context"
	"errors"
	"fmt"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
)

type ConfigReader interface {
	Load() (*config.Configuration, error)
}

func NewTrustCreator(commandIO *io.IO, configReader ConfigReader) *Creator {
	return &Creator{
		commandIO:    commandIO,
		configReader: configReader,
	}
}

type Creator struct {
	commandIO    *io.IO
	configReader ConfigReader
}

type CreatorOptions struct {
	IdentityFile string
}

func (t *Creator) CreateTrust(ctx context.Context, usernameRaw string, options CreatorOptions) error {
	_, err := t.loadConfig()
	if err != nil {
		return err
	}

	// Load local SSH key (verify?)
	// Get user ID
	// Look for standard crev-proofs repo
	// Present UI for rating
	// Present UI for comment
	// Sign
	// Write file
	// Commit

	return nil
}

func (t *Creator) loadConfig() (*config.Configuration, error) {
	conf, err := t.configReader.Load()
	if err != nil {
		return nil, err
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return conf, nil
}

func validateConfig(c *config.Configuration) error {
	// Should check if location exists in filesystem?
	if c.CurrentStore == "" {
		return errors.New("user current store is empty")
	}

	if c.CurrentID == nil {
		return errors.New("user current ID not set")
	}

	return nil
}
