package command

import (
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewStoreCommand() *cobra.Command {
	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "Manipulate proof stores",
	}

	storeCmd.AddCommand(NewFetchCommand())

	return storeCmd
}

const expectedFetchStoreArguments = 1

func NewFetchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch <url>",
		Short: "Fetch a proof store",
		RunE:  fetchStore,
		Args:  cobra.ExactArgs(expectedFetchStoreArguments),
	}
}

// args must be equal to length 1. This is ensured by cobra
func fetchStore(cmd *cobra.Command, args []string) error {
	fetcher := di.InitialiseStoreFetcher(io.New(cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr()))

	return fetcher.Fetch(cmd.Context(), args[0])
}
