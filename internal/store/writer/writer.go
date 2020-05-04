package writer

import "github.com/LaurenceGA/go-crev/proof/trust"

func New() *Writer {
	return &Writer{}
}

type Writer struct{}

func (w *Writer) SaveTrust(tr *trust.Trust) error {
	return nil
}
