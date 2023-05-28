package submission

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/cmd"
	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

// QueryOptions is an options to query a workspace.
type QueryOptions struct {
	WorkspaceName string
	RunID         string
	TaskName      string

	submissionClient factory.SubmissionClient
	workspaceClient  factory.WorkspaceClient
	formatter        formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewQueryOptions returns a reference to a QueryOptions
func NewQueryOptions(opt *clioptions.GlobalOptions) *QueryOptions {
	return &QueryOptions{
		options: opt,
	}
}

func NewCmdQuery(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewQueryOptions(opt)

	cmd := &cobra.Command{
		Use:   "query <submission_id>",
		Short: "get the status of the submission",
		Long:  "get the status of the submission",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.RunID, "run-id", "r", o.RunID, "The RunID of the submission.")
	cmd.Flags().StringVarP(&o.TaskName, "task-name", "t", o.TaskName, "The TaskName of the submission")

	return cmd
}

// Complete completes all the required options.
func (o *QueryOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workspaceClient, err = f.WorkspaceClient()
	if err != nil {
		return err
	}
	o.submissionClient, err = f.SubmissionClient()
	if err != nil {
		return err
	}

	if o.options.Stream.OutputFormat == "" {
		o.options.Stream.OutputFormat = o.GetDefaultFormat()
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the query options
func (o *QueryOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}

	if o.TaskName != "" {
		if o.RunID == "" {
			return fmt.Errorf("must specify a run id before specifying a task name ")
		}
	}

	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}

	return nil
}

// Run run the query workspace command
func (o *QueryOptions) Run(args []string) error {
	submissionID := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	resp, err := o.submissionClient.ListSubmissions(ctx, &convert.ListSubmissionsRequest{
		WorkspaceID: workspaceID,
		IDs:         []string{submissionID},
	})
	if err != nil {
		return err
	}

	if len(resp.Items) == 0 {
		return fmt.Errorf("submission: %s not found", submissionID)
	}

	if o.RunID != "" {
		runResp, err := o.submissionClient.ListRuns(ctx, &convert.ListRunsRequest{
			WorkspaceID:  workspaceID,
			SubmissionID: submissionID,
			IDs:          []string{o.RunID},
		})
		if err != nil {
			return err
		}

		if len(runResp.Items) == 0 {
			return fmt.Errorf("run: %s not found", o.TaskName)
		}

		if o.TaskName != "" {
			taskResp, err := o.submissionClient.ListTasks(ctx, &convert.ListTasksRequest{
				WorkspaceID:  workspaceID,
				SubmissionID: submissionID,
				RunID:        o.RunID,
			})
			if err != nil {
				return err
			}

			for _, item := range taskResp.Items {
				if item.Name == o.TaskName {
					o.formatter.Write(item.Status)
					return nil
				}
			}
			return fmt.Errorf("task: %s not found", o.TaskName)
		}

		o.formatter.Write(runResp.Items[0].Status)
		return nil
	}

	o.formatter.Write(resp.Items[0].Status)

	return nil
}

func (o *QueryOptions) GetPromptArgs() ([]string, error) {
	submissionID, err := prompt.PromptRequiredString("Submission ID")
	if err != nil {
		return []string{}, err
	}
	return []string{submissionID}, nil
}

func (o *QueryOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}
	o.RunID, err = prompt.PromptOptionalString("Run ID")
	if err != nil {
		return err
	}

	if o.RunID != "" {
		o.TaskName, err = prompt.PromptOptionalString("Task Name")
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *QueryOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
