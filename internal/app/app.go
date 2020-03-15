package app

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

func NewApplication(rootCmd *cobra.Command) *Application {
	return &Application{
		rootCmd: rootCmd,
	}
}

type Application struct {
	rootCmd *cobra.Command
}

func (a *Application) Execute() {
	err := a.rootCmd.ExecuteContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
