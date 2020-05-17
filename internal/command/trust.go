package command

import (
	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewTrustCommand() *cobra.Command {
	trustCmd := &cobra.Command{
		Use:   "trust",
		Short: "Manipulate trust proofs",
	}

	trustCmd.AddCommand(NewCreateTrustCommand())

	return trustCmd
}

const (
	expectedTrustArguments = 1

	identityFileFlagName = "identity-file"
)

func NewCreateTrustCommand() *cobra.Command {
	createTrustCmd := &cobra.Command{
		Use:   "create <Github username>",
		Short: "Create a trust proof for an ID",
		RunE:  createTrust,
		Args:  cobra.ExactArgs(expectedTrustArguments),
	}

	createTrustCmd.Flags().StringP(identityFileFlagName,
		"i",
		"",
		"identity file (private key) location to use for signing")

	return createTrustCmd
}

// args must have length equal to 1. This is ensured by cobra.
func createTrust(cmd *cobra.Command, args []string) error {
	trustCreator := di.InitialiseTrustCreator(ioFromCommand(cmd))

	idFilepath, err := cmd.Flags().GetString(identityFileFlagName)
	if err != nil {
		return err
	}

	return trustCreator.CreateTrust(
		cmd.Context(),
		args[0],
		trust.CreatorOptions{
			IdentityFile: idFilepath,
		})
}
