package command

import (
	"fmt"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/meta"
	"github.com/spf13/cobra"
)

const (
	// used by flows that sign proofs
	identityFileFlagName = "identity-file"
)

func NewRootCommand(commandIO *io.IO) *cobra.Command {
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
	rootCmd.AddCommand(NewIDCommand())
	rootCmd.AddCommand(NewVerifyCommand())
	rootCmd.AddCommand(NewTrustCommand())
	rootCmd.AddCommand(NewReviewCommand())

	rootCmd.SetIn(commandIO.In())
	rootCmd.SetOut(commandIO.Out())
	rootCmd.SetErr(commandIO.Err())

	rootCmd.PersistentFlags().BoolP(
		"verbose",
		"v",
		false,
		"show verbose output",
	)

	return rootCmd
}
