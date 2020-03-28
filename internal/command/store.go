package command

import (
	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewStoreCommand() *cobra.Command {
	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "Manipulate proof stores",
	}

	storeCmd.AddCommand(NewFetchCommand())
	storeCmd.AddCommand(NewSetCurrnetCommand())

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
	fetcher := di.InitialiseStoreFetcher(ioFromCommand(cmd))

	return fetcher.Fetch(cmd.Context(), args[0])
}

const expectedSetCurrentStoreArguments = 1

func NewSetCurrnetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-current <path>",
		Short: "Set the current user proof store",
		RunE:  setCurrentStore,
		Args:  cobra.ExactArgs(expectedSetCurrentStoreArguments),
	}
}

// args must be equal to length 1. This is ensured by cobra
func setCurrentStore(cmd *cobra.Command, args []string) error {
	configManipulator := di.InitialiseConfigManipulator()

	return configManipulator.SetCurrentStore(args[0])
}
