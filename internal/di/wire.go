//+build wireinject

// Package di contains all of the logic for how to wire up this application.
package di

import (
	"github.com/LaurenceGA/go-crev/internal/command/flow"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/store/fetcher"
	"github.com/LaurenceGA/go-crev/internal/verifier"
	"github.com/LaurenceGA/go-crev/internal/verifier/cloc"
	"github.com/LaurenceGA/go-crev/mod"
	"github.com/google/wire"
)

// InitialiseStoreFetcher create a fetcher for fetching crev proof stores
func InitialiseStoreFetcher(commandIO *io.IO) *fetcher.Fetcher {
	panic(wire.Build(fetcher.FetcherProvider))
}

func InitialiseVerifier(commandIO *io.IO) *verifier.Verifier {
	panic(wire.Build(
		verifier.New,

		wire.Bind(new(verifier.ModLister), new(*mod.Lister)),
		mod.NewLister,
		wire.Bind(new(mod.ModulesWrapper), new(*mod.Wrapper)),
		mod.NewWrapper,

		wire.Bind(new(verifier.LineCounter), new(*cloc.Counter)),
		cloc.New,
	))
}

func InitialiseConfigManipulator() *config.Manipulator {
	panic(wire.Build(
		config.ConfigManipulatorProvider,
	))
}

func InitialiseIDSetterFlow(commandIO *io.IO) *flow.IDSetter {
	panic(wire.Build(
		flow.NewIDSetter,

		wire.Bind(new(flow.ConfigManipulator), new(*config.Manipulator)),
		config.NewManipulator,

		wire.Bind(new(flow.Github), new(*github.Client)),
		github.NewClient,

		wire.Bind(new(flow.RepoFetcher), new(*fetcher.Fetcher)),
		fetcher.NewFetcher,

		wire.Bind(new(fetcher.GitCloner), new(*git.Client)),
		git.NewClient,

		wire.Bind(new(fetcher.FileDirs), new(*files.Filesystem)),  // Fetcher
		wire.Bind(new(config.FileFinder), new(*files.Filesystem)), // Config manipulator
		files.NewFilesystem,
		files.NewUserScope,
	))
}
