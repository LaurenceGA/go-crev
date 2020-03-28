// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/LaurenceGA/go-crev/internal/verifier"
	"github.com/LaurenceGA/go-crev/internal/verifier/cloc"
	"github.com/LaurenceGA/go-crev/mod"
)

// Injectors from wire.go:

func InitialiseStoreFetcher(commandIO *io.IO) *store.Fetcher {
	client := git.NewClient(commandIO)
	scope := files.NewUserScope()
	filesystem := files.NewFilesystem(scope)
	fetcher := store.NewFetcher(client, filesystem)
	return fetcher
}

func InitialiseVerifier(commandIO *io.IO) *verifier.Verifier {
	wrapper := mod.NewWrapper()
	lister := mod.NewLister(wrapper)
	counter := cloc.New()
	verifierVerifier := verifier.New(commandIO, lister, counter)
	return verifierVerifier
}

func InitialiseConfigManipulator() *config.Manipulator {
	scope := files.NewUserScope()
	filesystem := files.NewFilesystem(scope)
	manipulator := config.NewManipulator(filesystem)
	return manipulator
}
