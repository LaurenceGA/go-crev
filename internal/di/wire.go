//+build wireinject

// Package di contains all of the logic for how to wire up this application.
package di

import (
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/google/wire"
)

func InitialiseStoreFetcher() *store.Fetcher {
	panic(wire.Build(
		store.NewFetcher,

		wire.Bind(new(store.GitCloner), new(*git.Client)),
		git.NewClient,
	))
}
