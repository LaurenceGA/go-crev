package command

import (
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/spf13/cobra"
)

func ioFromCommand(cmd *cobra.Command) (*io.IO, error) {
	verbose, err := cmd.Root().PersistentFlags().GetBool(verboseOutputFlagName)
	if err != nil {
		return nil, err
	}

	return io.New(
		cmd.InOrStdin(),
		cmd.OutOrStdout(),
		cmd.ErrOrStderr(),
		verbose,
	), nil
}
