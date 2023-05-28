package data_model

import (
	"github.com/spf13/cobra"

	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

func NewCmdDataModel(opt *clioptions.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-model",
		Short: "data-model command",
		Long:  `data-model command`,
		Args:  cobra.NoArgs,
		Run:   prompt.SelectSubCommand,
	}
	cmd.AddCommand(NewCmdImport(opt))
	cmd.AddCommand(NewCmdList(opt))
	cmd.AddCommand(NewCmdDelete(opt))
	return cmd
}
