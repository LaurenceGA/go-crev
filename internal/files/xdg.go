// Package files directs the management of files owned by the application.
package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LaurenceGA/go-crev/meta"
	gap "github.com/muesli/go-app-paths"
)

type AppDirs interface {
	Data() (string, error)
	ConfigFile() (string, error)
}

func NewUserScope() *gap.Scope {
	return gap.NewScope(gap.User, meta.AppName)
}

func NewFilesystem(gapScope *gap.Scope) *Filesystem {
	return &Filesystem{
		scope: gapScope,
	}
}

type Filesystem struct {
	scope *gap.Scope
}

func (f *Filesystem) Data() (string, error) {
	dataDirs, err := f.scope.DataDirs()
	if err != nil {
		return "", err
	}

	// dataDirs is guaranteed to be of at least length 1
	return dataDirs[0], nil
}

const configFileName = "config.yml"

// ConfigFile will return the path to the highest precedence config file for the app.
// If no config file exists, one will be created in the default location.
func (f *Filesystem) ConfigFile() (string, error) {
	configs, err := f.scope.LookupConfig(configFileName)
	if err != nil {
		return "", fmt.Errorf("find valid config paths: %w", err)
	}

	if len(configs) > 0 {
		return configs[0], nil
	}

	newFilePath, err := f.scope.ConfigPath(configFileName)
	if err != nil {
		return "", fmt.Errorf("loading default config path: %w", err)
	}

	if err := createNewConfigFile(newFilePath); err != nil {
		return "", fmt.Errorf("creating new config file: %w", err)
	}

	return newFilePath, nil
}

// File owner can read write and executre the directory, others can't.
const ownerReadWrite = 0700

func createNewConfigFile(newFilePath string) error {
	if err := os.MkdirAll(filepath.Dir(newFilePath), ownerReadWrite); err != nil {
		return err
	}

	if _, err := os.Create(newFilePath); err != nil {
		return err
	}

	return nil
}
