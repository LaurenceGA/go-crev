package main

import (
	"log"

	"github.com/LaurenceGA/go-crev/internal/di"
)

func main() {
	app, err := di.InitialiseGoCrev()
	if err != nil {
		log.Fatal(err)
		return
	}

	app.Execute()
}
