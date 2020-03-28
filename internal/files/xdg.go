// Package files directs the management of files owned by the application
package files

import (
	"errors"

	"github.com/LaurenceGA/go-crev/meta"
	gap "github.com/muesli/go-app-paths"
)

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

	if len(dataDirs) == 0 {
		return "", errors.New("couldn't find a single data directory")
	}

	return dataDirs[0], nil
}
