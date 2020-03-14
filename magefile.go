// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/sh"
)

const thisRepo = "github.com/LaurenceGA/go-crev"

func Build() error {
	return sh.RunV("go", "build", "-ldflags", ldFlagsArg(), thisRepo+"/cmd/gocrev")
}

func ldFlagsArg() string {
	return constructLDFlags(getBuildDate(), "1.2.5")
}

func getBuildDate() string {
	return time.Now().Format(time.RFC3339)
}

func constructLDFlags(buildTime, version string) string {
	return fmt.Sprintf("-X '%s/version.BuildTime=%s' -X '%s/version.Version=%s'", thisRepo, buildTime, thisRepo, version)
}
