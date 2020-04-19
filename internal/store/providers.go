//+build wireinject

package store

import (
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/google/wire"
)

var FetcherProvider = wire.NewSet(
	NewFetcher,

	wire.Bind(new(GitCloner), new(*git.Client)),
	git.NewClient,

	wire.Bind(new(FileDirs), new(*files.Filesystem)),
	files.NewFilesystem,
	files.NewUserScope,
)
