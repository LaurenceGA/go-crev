package main

import (
	"context"
	"log"

	"github.com/LaurenceGA/go-crev/internal/command"
)

func main() {
	rootCmd := command.NewRootCommand(command.DefaultIO())

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
