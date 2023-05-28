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

// LogOptions is an options to log a workspace.
type LogOptions struct {
	WorkspaceName string
	RunID         string
	TaskName      string

	submissionClient factory.SubmissionClient
	workspaceClient  factory.WorkspaceClient
	formatter        formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewLogOptions returns a reference to a LogOptions
func NewLogOptions(opt *clioptions.GlobalOptions) *LogOptions {
	return &LogOptions{
		options: opt,
	}
}

func NewCmdLog(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewLogOptions(opt)

	cmd := &cobra.Command{
		Use:   "log <submission_id>",
		Short: "get the log of the submission",
		Long:  "get the log of the submission",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.RunID, "run-id", "r", o.RunID, "The RunID of the submission.")
	cmd.Flags().StringVarP(&o.TaskName, "task-name", "t", o.TaskName, "The TaskName of the submission")

	return cmd
}

// Complete completes all the required options.
func (o *LogOptions) Complete() error {
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

// Validate validate the log options
func (o *LogOptions) Validate() error {
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

// Run run the log workspace command
func (o *LogOptions) Run(args []string) error {
	submissionID := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
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
			return fmt.Errorf("run: %s not found", o.RunID)
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
					if item.Stdout != "" {
						o.formatter.Write(item.Stdout)
					} else {
						o.formatter.Write(fmt.Sprintf("the Stdout of task [%s] is unavailable now", o.TaskName))
					}
					if item.Stderr != "" {
						o.formatter.Write(item.Stderr)
					} else {
						o.formatter.Write(fmt.Sprintf("the Stderr of task [%s] is unavailable now", o.TaskName))
					}
					return nil
				}
			}

			if len(taskResp.Items) == 0 {
				return fmt.Errorf("task: %s not found", o.TaskName)
			}
			return nil
		}

		if runResp.Items[0].Log != nil {
			o.formatter.Write(runResp.Items[0].Log)
		} else {
			o.formatter.Write(fmt.Sprintf("the log of run [%s] is unavailable now", runResp.Items[0].ID))
		}
		return nil
	}

	runs, err := cmd.GetAllRuns(ctx, o.submissionClient, workspaceID, submissionID)
	if err != nil {
		return err
	}

	for _, run := range runs {
		if run.Log != nil {
			o.formatter.Write(struct {
				Log   *string
				RunID string
			}{
				Log:   run.Log,
				RunID: run.ID,
			})
			continue
		}
		o.formatter.Write(fmt.Sprintf("the log of run [%s] is unavailable now", run.ID))
	}

	return nil
}

func (o *LogOptions) GetPromptArgs() ([]string, error) {
	submissionID, err := prompt.PromptRequiredString("Submission ID")
	if err != nil {
		return []string{}, err
	}
	return []string{submissionID}, nil
}

func (o *LogOptions) GetPromptOptions() error {
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

func (o *LogOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
