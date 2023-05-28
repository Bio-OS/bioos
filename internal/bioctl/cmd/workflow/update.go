package workflow

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

// UpdateOptions is an options to update workflow.
type UpdateOptions struct {
	WorkspaceName string
	Name          string
	Description   string

	workflowClient  factory.WorkflowClient
	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

func (o *UpdateOptions) GetPromptArgs() ([]string, error) {
	workflowName, err := prompt.PromptRequiredString("Workflow", prompt.WithInputMessage("Name"))
	if err != nil {
		return nil, err
	}
	return []string{workflowName}, nil
}

func (o *UpdateOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = prompt.PromptRequiredString("Workspace", prompt.WithInputMessage("Name"))
	if err != nil {
		return err
	}
	o.Name, err = prompt.PromptRequiredString("Name")
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

// NewUpdateOptions returns a reference to a UpdateOptions.
func NewUpdateOptions(opt *clioptions.GlobalOptions) *UpdateOptions {
	return &UpdateOptions{
		options: opt,
	}
}

// NewCmdUpdate new a update workflow cmd.
func NewCmdUpdate(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewUpdateOptions(opt)

	cmd := &cobra.Command{
		Use:   "update <workflow_name>",
		Short: "update a workflow",
		Long:  "update a workflow",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The name of the workspace.")
	cmd.Flags().StringVarP(&o.Name, "name", "n", o.Name, "The name of the workflow.")
	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "The description of the workflow.")

	return cmd
}

// Complete completes all the required options.
func (o *UpdateOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workflowClient, err = f.WorkflowClient()
	if err != nil {
		return err
	}
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

// Run run the update workflow command
func (o *UpdateOptions) Run(args []string) error {
	workflowName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	workflowID, err := ConvertWorkflowNameIntoID(ctx, o.workflowClient, workspaceID, workflowName)
	if err != nil {
		return err
	}

	req := &convert.UpdateWorkflowRequest{
		ID:          workflowID,
		WorkspaceID: workspaceID,
	}
	if o.Name != "" {
		req.Name = o.Name
	}
	if o.Description != "" {
		req.Description = o.Description
	}
	// todo update api needs to add more fields
	_, err = o.workflowClient.UpdateWorkflow(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
