// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/LaurenceGA/go-crev/internal/app"
	"github.com/LaurenceGA/go-crev/internal/command"
)

// Injectors from wire.go:

func InitialiseGoCrev() (*app.Application, error) {
	commandIO := command.DefaultIO()
	cobraCommand := command.NewRootCommand(commandIO)
	application := app.NewApplication(cobraCommand)
	return application, nil
}
