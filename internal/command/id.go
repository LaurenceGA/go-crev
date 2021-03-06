package command

import (
	"fmt"

	"github.com/LaurenceGA/go-crev/internal/di"
	"github.com/spf13/cobra"
)

func NewIDCommand() *cobra.Command {
	idCmd := &cobra.Command{
		Use:   "id",
		Short: "Manipulate crev IDs",
	}

	idCmd.AddCommand(NewSetCurrentIDCommand())
	idCmd.AddCommand(NewShowCurrentIDCommand())
	idCmd.AddCommand(NewCreateNewIDCommand())

	return idCmd
}

func NewCreateNewIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Create a new CrevID",
		RunE:  createNewID,
	}
}

// args must have length equal to 1. This is ensured by cobra.
func createNewID(cmd *cobra.Command, args []string) error {
	panic("Implement me")
}

const expectedSetCurrentIDArguments = 1

func NewSetCurrentIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-current <Github username>",
		Short: "Set the current user ID",
		RunE:  setCurrentID,
		Args:  cobra.ExactArgs(expectedSetCurrentIDArguments),
	}
}

// args must have length equal to 1. This is ensured by cobra.
func setCurrentID(cmd *cobra.Command, args []string) error {
	commandIO, err := ioFromCommand(cmd)
	if err != nil {
		return err
	}

	setCurrentIDFlow := di.InitialiseIDSetterFlow(commandIO)

	return setCurrentIDFlow.SetFromUsername(cmd.Context(), args[0])
}

func NewShowCurrentIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the current user ID",
		RunE:  showCurrentID,
	}
}

func showCurrentID(cmd *cobra.Command, args []string) error {
	configManipulator := di.InitialiseConfigManipulator()

	curID, err := configManipulator.CurrentID()
	if err != nil {
		return err
	}

	if curID != nil {
		fmt.Fprintln(cmd.OutOrStdout(), "ID: "+curID.ID)

		if curID.Alias != "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Alias: "+curID.Alias)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Type: "+curID.Type)

		if curID.URL != "" {
			fmt.Fprintln(cmd.OutOrStdout(), "URL: "+curID.URL)
		}
	}

	return nil
}
