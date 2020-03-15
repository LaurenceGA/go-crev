package command

import (
	"fmt"

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

func NewFetchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch <url>",
		Short: "Fetch a proof store",
		RunE:  fetchStore,
	}
}

const expectedFetchStoreArguments = 1

func fetchStore(cmd *cobra.Command, args []string) error {
	if len(args) != expectedFetchStoreArguments {
		return fmt.Errorf("expected %d arguments, got %d", expectedFetchStoreArguments, len(args))
	}

	return nil
}
