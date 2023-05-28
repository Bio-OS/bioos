package prompt

import (
	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
)

const DefaultSelectSize = 10
const NeedPromptFlag = "needPrompt"

func SelectSubCommand(cmd *cobra.Command, args []string) {
	subCommands := cmd.Commands()
	selectedCmd, err := PromptSelectWithFunc("Select SubCommands:", DefaultSelectSize, func(i ...interface{}) ([]interface{}, error) {
		cmds := make([]interface{}, len(subCommands))
		for i, c := range subCommands {
			cmds[i] = c.Name()
		}
		return cmds, nil
	})
	utils.CheckErr(err)
	for i := range subCommands {
		if subCommands[i].Name() == selectedCmd {
			subCommands[i].Run(subCommands[i], []string{NeedPromptFlag})
		}
	}
}
