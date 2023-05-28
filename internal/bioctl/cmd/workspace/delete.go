package workspace

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/cmd"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
)

// DeleteOptions is an options to delete a workspace.
type DeleteOptions struct {
	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewDeleteOptions returns a reference to a DeleteOptions.
func NewDeleteOptions(opt *clioptions.GlobalOptions) *DeleteOptions {
	return &DeleteOptions{
		options: opt,
	}
}

// NewCmdDelete new a delete workspace cmd.
func NewCmdDelete(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewDeleteOptions(opt)

	cmd := &cobra.Command{
		Use:   "delete <workspace_name>",
		Short: "delete a workspace",
		Long:  "delete a workspace",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workspaceClient, err = f.WorkspaceClient()
	if err != nil {
		return err
	}

	if o.options.Stream.OutputFormat == "" {
		o.options.Stream.OutputFormat = o.GetDefaultFormat()
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the delete options
func (o *DeleteOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	return nil
}

// Run run the create workspace command
func (o *DeleteOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	workspaceID, err := ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, args[0])
	if err != nil {
		return err
	}

	_, err = o.workspaceClient.DeleteWorkspace(ctx, &convert.DeleteWorkspaceRequest{
		Id: workspaceID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) GetPromptArgs() ([]string, error) {
	name, err := cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return nil, err
	}

	return []string{name}, nil
}

func (o *DeleteOptions) GetPromptOptions() error {
	return nil
}

func (o *DeleteOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
