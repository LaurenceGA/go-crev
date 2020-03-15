//+build wireinject

package app

import (
	"github.com/LaurenceGA/go-crev/internal/command"
	"github.com/google/wire"
)

var (
	Provider = wire.NewSet(
		NewApplication,

		command.CommandProvider,
	)
)
