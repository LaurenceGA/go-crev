package command

import "github.com/spf13/cobra"

func InitialiseRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gocrev",
		Short: "A cryptographically verifiable code review system for go packages.",
		Long: `gocrev is a social code review system.

This is an implementation of the language agnostic crev.
For more information see here: https://github.com/crev-dev/crev

gocrev is in the early stages of development, please raise
issues here: https://github.com/LaurenceGA/go-crev/issues
`,
	}
}
