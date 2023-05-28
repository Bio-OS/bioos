package workspace

import (
	"github.com/spf13/cobra"

	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

func NewCmdWorkspace(opt *clioptions.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "workspace command",
		Long:  `workspace command`,
		Args:  cobra.NoArgs,
		Run:   prompt.SelectSubCommand,
	}
	cmd.AddCommand(NewCmdCreate(opt))
	cmd.AddCommand(NewCmdList(opt))
	cmd.AddCommand(NewCmdDelete(opt))
	cmd.AddCommand(NewCmdUpdate(opt))
	cmd.AddCommand(NewCmdImport(opt))
	return cmd
}
