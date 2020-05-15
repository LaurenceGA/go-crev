// Package cloc counts the lines of go source code.
package cloc

import (
	"fmt"

	"github.com/hhatto/gocloc"
)

func New() *Counter {
	return &Counter{}
}

type Counter struct {
}

type constError string

func (e constError) Error() string {
	return string(e)
}

const InvalidDirectory constError = "invalid directory"

// CountLines counts the lines of go source code in dir.
func (c *Counter) CountLines(dir string) (int, error) {
	const goLanguageKey = "Go"

	if dir == "" {
		return -1, InvalidDirectory
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
