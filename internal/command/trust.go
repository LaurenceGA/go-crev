package command

import (
	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

const (
	expectedTrustArguments = 1

	identityFileFlagName = "identity-file"
)

func NewTrustCommand() *cobra.Command {
	trustCmd := &cobra.Command{
		Use:   "trust <Github username>",
		Short: "Create a trust proof for an ID",
		RunE:  newTrust,
		Args:  cobra.ExactArgs(expectedTrustArguments),
	}

	trustCmd.Flags().StringP(identityFileFlagName,
		"i",
		"",
		"identity file (private key) location to use for signing")

	return trustCmd
}

// args must have length equal to 1. This is ensured by cobra
func newTrust(cmd *cobra.Command, args []string) error {
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
