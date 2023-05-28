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
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

// UpdateOptions is an options to update workspaces.
type UpdateOptions struct {
	Name        string
	Description string

	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewUpdateOptions returns a reference to a UpdateOptions.
func NewUpdateOptions(opt *clioptions.GlobalOptions) *UpdateOptions {
	return &UpdateOptions{
		options: opt,
	}
}

// NewCmdUpdate new a update workspace cmd.
func NewCmdUpdate(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewUpdateOptions(opt)

	cmd := &cobra.Command{
		Use:   "update <workspace_name>",
		Short: "update a workspace",
		Long:  "update a workspace",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.Name, "name", "n", o.Name, "The name of the workspace.")
	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "The description of the workspace.")

	return cmd
}

// Complete completes all the required options.
func (o *UpdateOptions) Complete() error {
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

// Validate validate the list options
func (o *UpdateOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	return nil
}

// Run run the update workspace command
func (o *UpdateOptions) Run(args []string) error {
	workspaceName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	workspaceID, err := ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, workspaceName)
	if err != nil {
		return err
	}

	req := &convert.UpdateWorkspaceRequest{
		ID: workspaceID,
	}
	if o.Name != "" {
		req.Name = &o.Name
	}
	if o.Description != "" {
		req.Description = &o.Description
	}

	_, err = o.workspaceClient.UpdateWorkspace(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (o *UpdateOptions) GetPromptArgs() ([]string, error) {
	name, err := cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return nil, err
	}

	return []string{name}, nil
}

func (o *UpdateOptions) GetPromptOptions() error {
	var err error

	o.Name, err = prompt.PromptOptionalString("Name")
	if err != nil {
		return err
	}

	o.Description, err = prompt.PromptOptionalString("Description")
	if err != nil {
		return err
	}

	return nil
}

func (o *UpdateOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
