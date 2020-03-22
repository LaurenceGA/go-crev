// Package cloc counts the lines of go source code
package cloc

import (
	"errors"
	"fmt"

	"github.com/hhatto/gocloc"
)

func New() *Counter {
	return &Counter{}
}

type Counter struct {
}

const goLanguageKey = "Go"

// CountLines counts the lines of go source code in dir
func (c *Counter) CountLines(dir string) (int, error) {
	if dir == "" {
		return -1, errors.New("invalid directory")
	}

	clocOpts := gocloc.NewClocOptions()
	clocOpts.IncludeLangs[goLanguageKey] = struct{}{}

	languages := gocloc.NewDefinedLanguages()
	processor := gocloc.NewProcessor(languages, clocOpts)

	result, err := processor.Analyze([]string{dir})
	if err != nil {
		return -1, fmt.Errorf("analysing '%s': %w", dir, err)
	}

	return int(result.Languages[goLanguageKey].Code), nil
}
