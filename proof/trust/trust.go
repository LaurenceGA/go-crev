package trust

import (
	"strings"
	"time"

	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/proof"
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
		proof.CommonData{
			Kind:    proof.Trust,
			Version: TrustVersion,
			Date:    time.Now(),
			From:    from,
		},
		Data{
			Level:   level,
			Comment: comment,
		},
	}
}

type Trust struct {
	proof.CommonData `yaml:",inline"`
	Data             `yaml:",inline"`
}

type Data struct {
	Level   Level
	Comment string
}

func (t *Trust) String() string {
	data, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(data)
}
