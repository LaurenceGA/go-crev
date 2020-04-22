//+build wireinject

package fetcher

import (
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/google/wire"
)

var FetcherProvider = wire.NewSet(
	NewFetcher,

	wire.Bind(new(GitCloner), new(*git.Client)),
	git.NewClient,

	wire.Bind(new(files.AppDirs), new(*files.Filesystem)),
	files.NewFilesystem,
	files.NewUserScope,
)
