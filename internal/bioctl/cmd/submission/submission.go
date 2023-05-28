package submission

import (
	"github.com/spf13/cobra"

	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

func NewCmdSubmission(opt *clioptions.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submission",
		Short: "submission command",
		Long:  `submission command`,
		Args:  cobra.NoArgs,
		Run:   prompt.SelectSubCommand,
	}
	cmd.AddCommand(NewCmdSubmit(opt))
	cmd.AddCommand(NewCmdQuery(opt))
	cmd.AddCommand(NewCmdDelete(opt))
	cmd.AddCommand(NewCmdStop(opt))
	cmd.AddCommand(NewCmdLog(opt))
	cmd.AddCommand(NewCmdList(opt))
	cmd.AddCommand(NewCmdOutput(opt))
	return cmd
}
