// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/LaurenceGA/go-crev/internal/command/flow/idset"
	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/store/fetcher"
	"github.com/LaurenceGA/go-crev/internal/verifier"
	"github.com/LaurenceGA/go-crev/internal/verifier/cloc"
	"github.com/LaurenceGA/go-crev/mod"
)

// Injectors from wire.go:

func InitialiseStoreFetcher(commandIO *io.IO) *fetcher.Fetcher {
	client := git.NewClient(commandIO)
	scope := files.NewUserScope()
	filesystem := files.NewFilesystem(scope)
	fetcherFetcher := fetcher.NewFetcher(client, filesystem)
	return fetcherFetcher
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

func InitialiseIDSetterFlow(commandIO *io.IO) *idset.IDSetter {
	scope := files.NewUserScope()
	filesystem := files.NewFilesystem(scope)
	manipulator := config.NewManipulator(filesystem)
	client := github.NewClient()
	gitClient := git.NewClient(commandIO)
	fetcherFetcher := fetcher.NewFetcher(gitClient, filesystem)
	idSetter := idset.NewIDSetter(commandIO, manipulator, client, fetcherFetcher)
	return idSetter
}

func InitialiseTrustCreator(commandIO *io.IO) *trust.Creator {
	scope := files.NewUserScope()
	filesystem := files.NewFilesystem(scope)
	manipulator := config.NewManipulator(filesystem)
	client := github.NewClient()
	creator := trust.NewTrustCreator(commandIO, manipulator, client)
	return creator
}
