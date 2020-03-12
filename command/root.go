package command

import (
	"github.com/LaurenceGA/go-crev/version"
	"github.com/spf13/cobra"
)

func InitialiseRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gocrev",
		Short: "A cryptographically verifiable code review system for go packages.",
		Long: `gocrev is a social code review system.

This is a go specific implementation of the language agnostic code review system crev.
For more information see here: https://github.com/crev-dev/crev

gocrev is in the early stages of development, please raise
issues here: https://github.com/LaurenceGA/go-crev/issues
`,
		Version: version.Version,
	}
}
