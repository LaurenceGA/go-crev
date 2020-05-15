package proof

import (
	"bytes"
	"fmt"
	"time"

	"github.com/LaurenceGA/go-crev/internal/id"
)

type Kind string

const (
	Trust  Kind = "trust"
	Review Kind = "package review"
)

type Proof interface {
	MarshalData() ([]byte, error)
	MarshalSignature() []byte
}

type CommonData struct {
	ID      string    `yaml:"id"`
	Kind    Kind      `yaml:"kind"`
	Version int       `yaml:"version"`
	Date    time.Time `yaml:"date"`
	From    id.ID     `yaml:"from"`
}

const (
	proofStart = "----- BEGIN CREV PROOF -----\n"
	proofSign  = "----- SIGN CREV PROOF -----\n"
	proofEnd   = "\n----- END CREV PROOF -----\n"
)

func MarhsalProof(p Proof) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(proofStart)

	data, err := p.MarshalData()
	if err != nil {
		return nil, fmt.Errorf("marshalling data: %w", err)
	}

	buf.Write(data)

	buf.WriteString(proofSign)

	sig := p.MarshalSignature()
	buf.Write(sig)

	buf.WriteString(proofEnd)

	return buf.Bytes(), nil
}
