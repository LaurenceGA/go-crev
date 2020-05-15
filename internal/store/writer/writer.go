package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/LaurenceGA/go-crev/proof"
	"github.com/LaurenceGA/go-crev/proof/trust"
)

const trustsPath = "trust"
const proofFileExtension = ".crev.proof"

func New() *Writer {
	return &Writer{}
}

type Writer struct{}

type constError string

func (e constError) Error() string {
	return string(e)
}

const NoValidIDOrAlias constError = "no valid ID or alias"

func (w *Writer) SaveTrust(dstStore *store.ProofStore, tr *trust.Trust) error {
	relPath, err := getRelativeTrustPath(tr)
	if err != nil {
		return fmt.Errorf("saving trust: %w", err)
	}

	trustPath := filepath.Join(dstStore.Dir, relPath)

	if err := os.MkdirAll(filepath.Dir(trustPath), 0755); err != nil {
		return err
	}

	marshalled, err := proof.MarhsalProof(tr)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(trustPath, marshalled, 0600); err != nil {
		return fmt.Errorf("writing trust file: %w", err)
	}

	// Commit

	return nil
}

type TrustFromIDTypeError struct {
	wanted, got id.CrevIdentity
}

func (e TrustFromIDTypeError) Error() string {
	return fmt.Sprintf("can't save trust of type %s, must be %s", e.got, e.wanted)
}

func getRelativeTrustPath(tr *trust.Trust) (string, error) {
	if tr.Data.From.Type != id.Github {
		return "", TrustFromIDTypeError{
			wanted: id.Github,
			got:    tr.Data.From.Type,
		}
	}

	idName, err := findIDName(&tr.Data.From)
	if err != nil {
		return "", fmt.Errorf("creating trust path: %w", err)
	}

	return filepath.Join(idName, trustsPath, yearMonthTimestamp(), tr.Data.ID+proofFileExtension), nil
}

func yearMonthTimestamp() string {
	return time.Now().Format("2006-01")
}

func findIDName(i *id.ID) (string, error) {
	if i.Alias != "" {
		return i.Alias, nil
	}

	if i.ID != "" {
		return i.ID, nil
	}

	return "", NoValidIDOrAlias
}
