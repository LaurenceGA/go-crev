//+build wireinject

package command

import (
	"github.com/google/wire"
)

var (
	CommandProvider = wire.NewSet(
		NewRootCommand,

		DefaultIO,
	)
)
