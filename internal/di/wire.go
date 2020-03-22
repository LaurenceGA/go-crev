//+build wireinject

// Package di contains all of the logic for how to wire up this application.
package di

import (
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/mod"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/LaurenceGA/go-crev/internal/verifier"
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
	))
}

func InitialiseVerifier(commandIO *io.IO) *verifier.Verifier {
	panic(wire.Build(
		verifier.New,

		wire.Bind(new(verifier.ModLister), new(*mod.Lister)),
		mod.NewLister,
	))
}
