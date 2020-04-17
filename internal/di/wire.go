//+build wireinject

// Package di contains all of the logic for how to wire up this application.
package di

import (
	"github.com/LaurenceGA/go-crev/internal/command/flow"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/LaurenceGA/go-crev/internal/verifier"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/verifier/cloc"
	"github.com/LaurenceGA/go-crev/mod"
	"github.com/google/wire"
)

// InitialiseStoreFetcher create a fetcher for fetching crev proof stores
func InitialiseStoreFetcher(commandIO *io.IO) *store.Fetcher {
	panic(wire.Build(
		store.NewFetcher,

		wire.Bind(new(store.GitCloner), new(*git.Client)),
		git.NewClient,

		wire.Bind(new(store.FileDirs), new(*files.Filesystem)),
		files.NewFilesystem,
		files.NewUserScope,
	))
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
	panic(wire.Build(config.ConfigManipulatorProvider))
}

func InitialiseIDSetterFlow() *flow.IDSetter {
	panic(wire.Build(
		flow.NewIDSetter,

		wire.Bind(new(flow.ConfigManipulator), new(*config.Manipulator)),
		config.ConfigManipulatorProvider,

		wire.Bind(new(flow.GithubUser), new(*github.Client)),
		github.NewClient,
	))
}
