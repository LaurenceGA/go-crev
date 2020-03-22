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
	verifier := di.InitialiseVerifier(ioFromCommand(cmd))

	return verifier.Verify()
}
