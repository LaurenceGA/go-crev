//+build wireinject

package config

import (
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/google/wire"
)

var ConfigManipulatorProvider = wire.NewSet(
	NewManipulator,

	wire.Bind(new(files.AppDirs), new(*files.Filesystem)),
	files.NewFilesystem,
	files.NewUserScope,
)
