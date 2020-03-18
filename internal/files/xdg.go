// Package files directs the management of files owned by the application
package files

import "github.com/adrg/xdg"

func NewFilesystem() *Filesystem {
	return &Filesystem{}
}

type Filesystem struct {
}

func (f *Filesystem) Cache() string {
	return xdg.CacheHome
}
