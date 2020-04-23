package trust

import (
	"context"

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
	_, err := t.configReader.Load()
	if err != nil {
		return err
	}

	// Load local ID
	// Load local store location
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
