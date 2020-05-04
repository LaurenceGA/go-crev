package writer

import (
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/LaurenceGA/go-crev/proof/trust"
)

func New() *Writer {
	return &Writer{}
}

type Writer struct{}

func (w *Writer) SaveTrust(dstStore *store.ProofStore, tr *trust.Trust) error {
	// Serialize
	// Commit
	return nil
}
