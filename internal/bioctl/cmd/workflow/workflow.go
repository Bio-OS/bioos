package workflow

import (
	"github.com/spf13/cobra"

	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

func NewCmdWorkflow(opt *clioptions.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "workflow command",
		Long:  `workflow command`,
		Args:  cobra.NoArgs,
		Run:   prompt.SelectSubCommand,
	}
	cmd.AddCommand(NewCmdCreate(opt))
	cmd.AddCommand(NewCmdImport(opt))
	cmd.AddCommand(NewCmdList(opt))
	cmd.AddCommand(NewCmdUpdate(opt))
	cmd.AddCommand(NewCmdDelete(opt))
	return cmd
}
