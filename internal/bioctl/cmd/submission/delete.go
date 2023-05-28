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

// DeleteOptions is an options to delete a workspace.
type DeleteOptions struct {
	WorkspaceName string

	submissionClient factory.SubmissionClient
	workspaceClient  factory.WorkspaceClient
	formatter        formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewDeleteOptions returns a reference to a DeleteOptions
func NewDeleteOptions(opt *clioptions.GlobalOptions) *DeleteOptions {
	return &DeleteOptions{
		options: opt,
	}
}

func NewCmdDelete(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewDeleteOptions(opt)

	cmd := &cobra.Command{
		Use:   "delete <submission_id>",
		Short: "delete the submission",
		Long:  "delete the submission",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")

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

// Validate validate the delete options
func (o *DeleteOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}

	return nil
}

// Run run the delete workspace command
func (o *DeleteOptions) Run(args []string) error {
	submissionID := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	_, err = o.submissionClient.DeleteSubmission(ctx, &convert.DeleteSubmissionRequest{
		WorkspaceID: workspaceID,
		ID:          submissionID,
	})
	if err != nil {
		return err
	}

	o.formatter.Write(fmt.Sprintf("submission [%s] will be deleted soon", submissionID))

	return nil
}

func (o *DeleteOptions) GetPromptArgs() ([]string, error) {
	submissionID, err := prompt.PromptRequiredString("Submission ID")
	if err != nil {
		return []string{}, err
	}
	return []string{submissionID}, nil
}

func (o *DeleteOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}
	return nil
}

func (o *DeleteOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
