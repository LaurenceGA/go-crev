package command

import (
	"fmt"

	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewStoreCommand() *cobra.Command {
	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "Manipulate proof stores",
	}

	storeCmd.AddCommand(NewFetchCommand())
	storeCmd.AddCommand(NewSetCurrentStoreCommand())
	storeCmd.AddCommand(NewShowCurrentStoreCommand())

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

// args must have length equal to 1. This is ensured by cobra.
func fetchStore(cmd *cobra.Command, args []string) error {
	fetcher := di.InitialiseStoreFetcher(ioFromCommand(cmd))

	_, err := fetcher.Fetch(cmd.Context(), args[0])

	return err
}

const expectedSetCurrentStoreArguments = 1

func NewSetCurrentStoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-current <path>",
		Short: "Set the current user proof store",
		RunE:  setCurrentStore,
		Args:  cobra.ExactArgs(expectedSetCurrentStoreArguments),
	}
}

// args must have length equal to 1. This is ensured by cobra.
func setCurrentStore(cmd *cobra.Command, args []string) error {
	configManipulator := di.InitialiseConfigManipulator()

	return configManipulator.SetCurrentStore(args[0])
}

func NewShowCurrentStoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the current user proof store",
		RunE:  showCurrentStore,
	}
}

func showCurrentStore(cmd *cobra.Command, args []string) error {
	configManipulator := di.InitialiseConfigManipulator()

	curStore, err := configManipulator.CurrentStore()
	if err != nil {
		return err
	}

	if curStore != "" {
		fmt.Fprintln(cmd.OutOrStdout(), curStore)
	}

	return nil
}
