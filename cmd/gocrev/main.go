package main

import (
	"context"
	"log"

	"github.com/LaurenceGA/go-crev/internal/command"
	"github.com/LaurenceGA/go-crev/internal/command/io"
)

func main() {
	rootCmd := command.NewRootCommand(&io.IO{})

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
