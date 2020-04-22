package flow

import (
	"context"

	"github.com/LaurenceGA/go-crev/internal/command/io"
)

func NewTrustCreator(commandIO *io.IO) *TrustCreator {
	return &TrustCreator{
		commandIO: commandIO,
	}
}

type TrustCreator struct {
	commandIO *io.IO
}

func (t *TrustCreator) CreateTrust(ctx context.Context, usernameRaw string) error {
	return nil
}
