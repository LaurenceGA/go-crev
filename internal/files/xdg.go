// Package files directs the management of files owned by the application
package files

import (
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

func (f *Filesystem) Cache() (string, error) {
	return f.scope.CacheDir()
}
