package command

import (
	"fmt"
	"io"
	"os"

	"github.com/LaurenceGA/go-crev/version"
	"github.com/spf13/cobra"
)

func DefaultIO() *CommandIO {
	return &CommandIO{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
}

type CommandIO struct {
	in       io.Reader
	out, err io.Writer
}

func NewRootCommand(commandIO *CommandIO) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gocrev",
		Short: "A cryptographically verifiable code review system for go packages.",
		Long: fmt.Sprintf(`gocrev is a social code review system.

This is a go specific implementation of the language agnostic code review system crev.
For more information see here: https://github.com/crev-dev/crev

gocrev is in the early stages of development, please raise
issues here: https://github.com/LaurenceGA/go-crev/issues

This is version %s, built at %s
`, version.Version, version.BuildTime),
		Version: version.Version,
	}

	rootCmd.AddCommand(NewStoreCommand())
	rootCmd.SetIn(commandIO.in)
	rootCmd.SetOut(commandIO.out)
	rootCmd.SetErr(commandIO.err)

	return rootCmd
}
