package trust

import (
	"context"

	"github.com/LaurenceGA/go-crev/internal/command/io"
)

func NewTrustCreator(commandIO *io.IO) *Creator {
	return &Creator{
		commandIO: commandIO,
	}
}

type Creator struct {
	commandIO *io.IO
}

type CreatorOptions struct {
	IdentityFile string
}

func (t *Creator) CreateTrust(ctx context.Context, usernameRaw string, options CreatorOptions) error {
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
