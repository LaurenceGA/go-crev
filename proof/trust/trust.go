package trust

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/proof"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

const TrustVersion = -1

type Level string

const (
	Distrust Level = "distrust"
	None     Level = "none"
	Low      Level = "low"
	Medium   Level = "medium"
	High     Level = "high"
)

func levelLookupMap() map[string]Level {
	levels := Levels()
	lookupMap := make(map[string]Level)

	for _, l := range levels {
		lookupMap[string(l)] = l
	}

	return lookupMap
}

func Levels() []Level {
	return []Level{Distrust, None, Low, Medium, High}
}

func ToLevel(s string) (Level, bool) {
	l, ok := levelLookupMap()[strings.ToLower(s)]

	return l, ok
}

func New(from id.ID, level Level, comment string) *Trust {
	return &Trust{
		Data: Data{
			Level:   level,
			Comment: comment,
			CommonData: proof.CommonData{
				Kind:    proof.Trust,
				Version: TrustVersion,
				Date:    time.Now(),
				From:    from,
			},
		},
	}
}

type Trust struct {
	Data      Data
	Signature string
}

type Data struct {
	Level            Level
	Comment          string
	proof.CommonData `yaml:",inline"`
}

func (t *Trust) Sign(signer ssh.Signer) error {
	data, err := yaml.Marshal(t.Data)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}

	signature, err := signer.Sign(rand.Reader, data)
	if err != nil {
		return fmt.Errorf("signing message: %w", err)
	}

	t.Signature = base64.StdEncoding.EncodeToString(signature.Blob)

	return nil
}

func (t *Trust) String() string {
	data, err := yaml.Marshal(t.Data)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s\n%s", string(data), t.Signature)
}
