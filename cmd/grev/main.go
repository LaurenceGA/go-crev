package main

import (
	"log"

	"github.com/LaurenceGA/go-crev/command"
)

func main() {
	rootCmd := command.InitialiseRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
