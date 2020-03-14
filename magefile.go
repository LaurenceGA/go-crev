// +build mage

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/magefile/mage/sh"
)

const thisRepo = "github.com/LaurenceGA/go-crev"

func Build() error {
	return sh.RunV("go", "build", "-ldflags", ldFlagsArg(), thisRepo+"/cmd/gocrev")
}

func ldFlagsArg() string {
	buildTime := getBuildTime()
	version := getVersion()

	log.Printf("Setting builtTime to %s\n", buildTime)
	log.Printf("Setting version to %s\n", version)

	return constructLDFlags(buildTime, version)
}

func getBuildTime() string {
	return time.Now().Format(time.RFC3339)
}

const gitNoTagsStatus = 128

func getVersion() string {
	latestTag, err := getLatestTag()
	if err != nil {
		if sh.ExitStatus(errors.Unwrap(err)) == gitNoTagsStatus {
			log.Println("No tags found")

			return getShortHash()
		} else {
			log.Fatal(err)
		}
	}

	return fmt.Sprintf("%s-%s", latestTag, getShortHash())
}

// gets the latest tag for this branch
func getLatestTag() (string, error) {
	var buf bytes.Buffer
	_, err := sh.Exec(nil, &buf, ioutil.Discard, "git", "describe", "--abbrev=0", "--tags")
	if err != nil {
		return "", fmt.Errorf("retrieving tag: %w", err)
	}

	return buf.String(), nil
}

func getShortHash() string {
	output, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

	return output
}

func constructLDFlags(buildTime, version string) string {
	return fmt.Sprintf("-X '%s/version.BuildTime=%s' -X '%s/version.Version=%s'", thisRepo, buildTime, thisRepo, version)
}
