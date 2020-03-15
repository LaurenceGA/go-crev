//+build wireinject

// Package di contains all of the logic for how to wire up this application.
package di

import (
	"github.com/LaurenceGA/go-crev/internal/app"
	"github.com/google/wire"
)

// InitialiseGoCrev creates an initialised go crev application.
func InitialiseGoCrev() (*app.Application, error) {
	panic(wire.Build(app.Provider))
}
