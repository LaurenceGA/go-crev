package command

import (
	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify the trustworthiness of this project's dependencies.",
		RunE:  verifyCurrentPackage,
	}

	return verifyCmd
}

func verifyCurrentPackage(cmd *cobra.Command, args []string) error {
	commandIO, err := ioFromCommand(cmd)
	if err != nil {
		return err
	}

	verifier := di.InitialiseVerifier(commandIO)

	return verifier.Verify()
}
