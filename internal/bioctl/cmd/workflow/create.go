package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
)

// CreateOptions is an options to create a workflow.
type CreateOptions struct {
	WorkspaceName    string
	Name             string
	Description      string
	Language         string
	Source           string
	URL              string
	Tag              string
	Token            string
	MainWorkflowPath string

	workflowClient  factory.WorkflowClient
	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

func (o *CreateOptions) GetPromptArgs() ([]string, error) {
	workflowName, err := prompt.PromptRequiredString("Name")
	if err != nil {
		return []string{}, err
	}
	return []string{workflowName}, nil
}

func (o *CreateOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = prompt.PromptRequiredString("Workspace", prompt.WithInputMessage("Name"))
	if err != nil {
		return err
	}
	o.Description, err = prompt.PromptOptionalString("Description")
	if err != nil {
		return err
	}

	o.Language, err = prompt.PromptStringSelect("Language", 10, []string{workflow.LanguageWDL, workflow.LanguageNextflow, workflow.LanguageCWL, workflow.LanguageSnakemake})
	if err != nil {
		return err
	}

	o.Source, err = prompt.PromptStringSelect("Source", 2, []string{workflow.WorkflowSourceGit, workflow.WorkflowSourceFile})
	if err != nil {
		return err
	}
	o.URL, err = prompt.PromptRequiredString("URL")
	if err != nil {
		return err
	}
	o.Tag, err = prompt.PromptRequiredString("Tag", prompt.WithInputMessage("tag or branch"))
	if err != nil {
		return err
	}
	o.Token, err = prompt.PromptOptionalString("Token")
	if err != nil {
		return err
	}
	o.MainWorkflowPath, err = prompt.PromptRequiredString("MainWorkflowPath", prompt.WithInputMessage("relative path"))
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}

// NewCreateOptions returns a reference to a CreateOptions
func NewCreateOptions(opt *clioptions.GlobalOptions) *CreateOptions {
	return &CreateOptions{
		options: opt,
	}
}

func NewCmdCreate(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewCreateOptions(opt)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a workflow",
		Long:  "create a workflow",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The name of the workspace.")
	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "The description of the workflow.")
	cmd.Flags().StringVarP(&o.Language, "language", "l", o.Language, "The language of the workflow.")
	cmd.Flags().StringVar(&o.Source, "source", o.Source, "The source of the workflow.")
	cmd.Flags().StringVar(&o.URL, "url", o.URL, "The url of the workflow.")
	cmd.Flags().StringVarP(&o.Tag, "tag", o.Tag, "t", "The tag or branch of the workflow.")
	cmd.Flags().StringVar(&o.Token, "token", o.Token, "The token to clone repo of the workflow.")
	cmd.Flags().StringVarP(&o.MainWorkflowPath, "main-workflow-path", "p", o.MainWorkflowPath, "The main path of the workflow.")

	return cmd
}

// Complete completes all the required options.
func (o *CreateOptions) Complete() error {
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

// Validate validate the create options
func (o *CreateOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.Language != workflow.LanguageWDL || o.Language != workflow.LanguageNextflow || o.Language != workflow.LanguageCWL || o.Language != workflow.LanguageSnakemake {
		return fmt.Errorf("unspport language: %s", o.Language)
	}
	if o.Source != workflow.WorkflowSourceGit && o.Source != workflow.WorkflowSourceFile {
		return fmt.Errorf("unspport source: %s", o.Source)
	}
	return nil
}

// Run run the create workflow command
func (o *CreateOptions) Run(args []string) error {
	workflowName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}
	req := &convert.CreateWorkflowRequest{
		WorkspaceID:      workspaceID,
		Name:             workflowName,
		Description:      o.Description,
		Language:         o.Language,
		Source:           o.Source,
		MainWorkflowPath: o.MainWorkflowPath,
		Url:              o.URL,
		Tag:              o.Tag,
		Token:            o.Token,
	}

	var resp *convert.CreateWorkflowResponse
	resp, err = o.workflowClient.CreateWorkflow(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp.ID)

	return nil
}
