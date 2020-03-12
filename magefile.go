// +build mage

package main

import "github.com/magefile/mage/sh"

func Build() error {
	return sh.RunV("go", "build", "github.com/LaurenceGA/go-crev/cmd/gocrev")
}