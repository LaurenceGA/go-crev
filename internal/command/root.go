package command

import (
	"fmt"
	"io"
	"os"

	"github.com/LaurenceGA/go-crev/meta"
	"github.com/spf13/cobra"
)

func DefaultIO() *IO {
	return &IO{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
}

type IO struct {
	in       io.Reader
	out, err io.Writer
}

func NewRootCommand(commandIO *IO) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   meta.AppName,
		Short: "A cryptographically verifiable code review system for go packages.",
		Long: fmt.Sprintf(`%s is a social code review system.

This is a go specific implementation of the language agnostic code review system crev.
For more information see here: https://github.com/crev-dev/crev

%s is in the early stages of development, please raise
issues here: https://github.com/LaurenceGA/go-crev/issues

This is version %s, built at %s
`, meta.AppName, meta.AppName, meta.Version, meta.BuildTime),
		Version: meta.Version,
	}

	rootCmd.AddCommand(NewStoreCommand())
	rootCmd.SetIn(commandIO.in)
	rootCmd.SetOut(commandIO.out)
	rootCmd.SetErr(commandIO.err)

	return rootCmd
}
