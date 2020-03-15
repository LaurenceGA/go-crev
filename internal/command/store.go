package command

import (
	"fmt"

	"github.com/LaurenceGA/go-crev/internal/store"
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

type Fetcher interface {
	Fetch(string) error
}

type FetchStoreCommand struct {
	fetcher Fetcher
}

func NewFetchCommand() *cobra.Command {
	cmd := &FetchStoreCommand{
		fetcher: &store.Fetcher{},
	}

	return &cobra.Command{
		Use:   "fetch <url>",
		Short: "Fetch a proof store",
		RunE:  cmd.fetchStore,
	}
}

const expectedFetchStoreArguments = 1

func (f *FetchStoreCommand) fetchStore(cmd *cobra.Command, args []string) error {
	if len(args) != expectedFetchStoreArguments {
		return fmt.Errorf("expected %d arguments, got %d", expectedFetchStoreArguments, len(args))
	}

	return f.fetcher.Fetch(args[0])
}
