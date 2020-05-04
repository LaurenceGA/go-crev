package proof

import (
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
	Signature() string
}

type CommonData struct {
	Kind    Kind      `yaml:"kind"`
	Version int       `yaml:"version"`
	Date    time.Time `yaml:"date"`
	From    id.ID     `yaml:"from"`
}
