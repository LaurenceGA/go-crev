package command

import (
	"github.com/spf13/cobra"
)

const expectedTrustArguments = 1

func NewTrustCommand() *cobra.Command {
	trustCmd := &cobra.Command{
		Use:   "trust <Github username>",
		Short: "Create a trust proof for an ID",
		RunE:  newTrust,
		Args:  cobra.ExactArgs(expectedTrustArguments),
	}

	return trustCmd
}

// args must have length equal to 1. This is ensured by cobra
func newTrust(cmd *cobra.Command, args []string) error {
	return nil
}
