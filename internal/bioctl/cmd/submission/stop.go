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

// StopOptions is an options to stop a workspace.
type StopOptions struct {
	WorkspaceName string
	RunID         string

	submissionClient factory.SubmissionClient
	workspaceClient  factory.WorkspaceClient
	formatter        formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewStopOptions returns a reference to a StopOptions
func NewStopOptions(opt *clioptions.GlobalOptions) *StopOptions {
	return &StopOptions{
		options: opt,
	}
}

func NewCmdStop(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewStopOptions(opt)

	cmd := &cobra.Command{
		Use:   "stop <submission_id>",
		Short: "stop the submission",
		Long:  "stop the submission",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.RunID, "run-id", "r", o.RunID, "The RunID of the submission.")

	return cmd
}

// Complete completes all the required options.
func (o *StopOptions) Complete() error {
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

// Validate validate the stop options
func (o *StopOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}

	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}
	return nil
}

// Run run the stop workspace command
func (o *StopOptions) Run(args []string) error {
	submissionID := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	if o.RunID != "" {
		_, err := o.submissionClient.CancelRun(ctx, &convert.CancelRunRequest{
			WorkspaceID:  workspaceID,
			SubmissionID: submissionID,
			ID:           o.RunID,
		})
		if err != nil {
			return err
		}

		o.formatter.Write(fmt.Sprintf("run [%s] will be canceled soon", o.RunID))

		return nil
	}

	_, err = o.submissionClient.CancelSubmission(ctx, &convert.CancelSubmissionRequest{
		WorkspaceID: workspaceID,
		ID:          submissionID,
	})
	if err != nil {
		return err
	}

	o.formatter.Write(fmt.Sprintf("submission [%s] will be canceled soon", submissionID))

	return nil
}

func (o *StopOptions) GetPromptArgs() ([]string, error) {
	submissionID, err := prompt.PromptRequiredString("Submission ID")
	if err != nil {
		return []string{}, err
	}
	return []string{submissionID}, nil
}

func (o *StopOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}
	o.RunID, err = prompt.PromptOptionalString("Run ID")
	if err != nil {
		return err
	}
	return nil
}

func (o *StopOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
