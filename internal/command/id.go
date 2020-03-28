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

	return idCmd
}

const expectedSetCurrentIDArguments = 1

func NewSetCurrentIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-current <path>",
		Short: "Set the current user ID",
		RunE:  setCurrentID,
		Args:  cobra.ExactArgs(expectedSetCurrentIDArguments),
	}
}

// args must be equal to length 1. This is ensured by cobra
func setCurrentID(cmd *cobra.Command, args []string) error {
	configManipulator := di.InitialiseConfigManipulator()

	return configManipulator.SetCurrentID(args[0])
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

	if curID != "" {
		fmt.Fprintln(cmd.OutOrStdout(), curID)
	}

	return nil
}
