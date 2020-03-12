package main

import (
	"context"
	"log"

	"github.com/LaurenceGA/go-crev/command"
)

func main() {
	rootCmd := command.InitialiseRootCommand()

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}