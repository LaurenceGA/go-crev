package command

import (
	"github.com/LaurenceGA/go-crev/internal/command/flow/review"
	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewReviewCommand() *cobra.Command {
	reviewCmd := &cobra.Command{
		Use:   "review",
		Short: "Manipulate package review proofs",
	}

	reviewCmd.AddCommand(NewCreateReviewCommand())

	return reviewCmd
}

const expectedReviewArguments = 1

func NewCreateReviewCommand() *cobra.Command {
	createReviewCmd := &cobra.Command{
		Use:   "create <package>",
		Short: "Create a package reivew proof for a package",
		RunE:  createReview,
		Args:  cobra.ExactArgs(expectedReviewArguments),
	}

	createReviewCmd.Flags().StringP(identityFileFlagName,
		"i",
		"",
		"identity file (private key) location to use for signing")

	return createReviewCmd
}

// args must have length equal to 1. This is ensured by cobra.
func createReview(cmd *cobra.Command, args []string) error {
	reviewCreator := di.InitialiseReviewCreator(ioFromCommand(cmd))

	idFilepath, err := cmd.Flags().GetString(identityFileFlagName)
	if err != nil {
		return err
	}

	_ = idFilepath

	return reviewCreator.CreateReview(
		cmd.Context(),
		args[0],
		review.CreatorOptions{
			IdentityFile: idFilepath,
		},
	)
}
