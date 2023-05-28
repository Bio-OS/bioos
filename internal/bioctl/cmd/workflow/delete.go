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

// DeleteOptions is an options to delete a workflow.
type DeleteOptions struct {
	WorkspaceName string

	workflowClient  factory.WorkflowClient
	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

func (o *DeleteOptions) GetPromptArgs() ([]string, error) {
	workflowName, err := prompt.PromptRequiredString("Workflow", prompt.WithInputMessage("Name"))
	if err != nil {
		return []string{}, err
	}
	return []string{workflowName}, nil
}

func (o *DeleteOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = prompt.PromptRequiredString("Workspace", prompt.WithInputMessage("Name"))

	return err
}

func (o *DeleteOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}

// NewDeleteOptions returns a reference to a DeleteOptions.
func NewDeleteOptions(opt *clioptions.GlobalOptions) *DeleteOptions {
	return &DeleteOptions{
		options: opt,
	}
}

// NewCmdDelete new a delete workflow cmd.
func NewCmdDelete(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewDeleteOptions(opt)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a workflow",
		Long:  "delete a workflow",
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The name of the workspace.")

	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete() error {
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

// Validate validate the delete options
func (o *DeleteOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	return nil
}

// Run run the create workflow command
func (o *DeleteOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}
	if len(args) != 0 {
		workflowName := args[0]
		var workflowID string
		workflowID, err = ConvertWorkflowNameIntoID(ctx, o.workflowClient, workspaceID, workflowName)
		if err != nil {
			return err
		}
		_, err = o.workflowClient.DeleteWorkflow(ctx, &convert.DeleteWorkflowRequest{
			ID:          workflowID,
			WorkspaceID: workspaceID,
		})
	} else {
		// delete all workflow in special workspace
		pageNum := 1
		pageSize := 100
		for {
			var workspaces *convert.ListWorkflowsResponse
			workspaces, err = o.workflowClient.ListWorkflow(ctx, &convert.ListWorkflowsRequest{
				Page:        pageNum,
				Size:        pageSize,
				WorkspaceID: workspaceID,
			})
			if err != nil {
				return err
			}
			for _, item := range workspaces.Items {
				_, err = o.workflowClient.DeleteWorkflow(ctx, &convert.DeleteWorkflowRequest{
					ID:          item.ID,
					WorkspaceID: workspaceID,
				})
			}
			if len(workspaces.Items) < pageSize {
				return nil
			}
			pageNum++
		}
	}

	return err
}
