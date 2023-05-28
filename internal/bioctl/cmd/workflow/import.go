package workflow

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

// ImportOptions is an options to get a workflow.
type ImportOptions struct {
	WorkspaceName string
	Filepath      string

	create *CreateOptions

	workflowClient  factory.WorkflowClient
	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

func (o *ImportOptions) GetPromptArgs() ([]string, error) {
	return []string{}, nil
}

func (o *ImportOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = prompt.PromptRequiredString("Workspace", prompt.WithInputMessage("Name"))
	if err != nil {
		return err
	}
	o.Filepath, err = prompt.PromptStringWithValidator("Input File Path", func(ans interface{}) error {
		err := survey.Required(ans)
		if err != nil {
			return err
		}
		_, err = os.Stat(cast.ToString(ans))
		if err != nil {
			curPath, _ := os.Getwd()
			return fmt.Errorf("%w, (currenct path is [%s])", err, curPath)
		}
		return nil
	}, prompt.WithInputMessage("abs path"))

	if err != nil {
		return err
	}

	return nil
}

func (o *ImportOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}

// NewImportOptions returns a reference to a ImportOptions.
func NewImportOptions(opt *clioptions.GlobalOptions) *ImportOptions {
	return &ImportOptions{
		options: opt,
	}
}

// NewCmdImport new a get workflow cmd.
func NewCmdImport(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewImportOptions(opt)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "import a workflow",
		Long:  "import a workflow",
		Example: `# import a workflow
		# import a workflow with json format
		{"name":"test","description":"test","language":"WDL","source":"git","URL":"https://github.com/xueerli/wdl.git","tag":"main","token":"","mainWorkflowPath":"no.wdl"}
`,
		Args: cobra.ExactArgs(0),
		Run:  clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.Filepath, "file", "f", o.WorkspaceName, "The file of the git info.")
	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The name of workspace")

	return cmd
}

// Complete completes all the required options.
func (o *ImportOptions) Complete() error {
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
	if o.options.Stream.OutputFormat == "" {
		o.options.Stream.OutputFormat = o.GetDefaultFormat()
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the get options
func (o *ImportOptions) Validate() error {
	stat, err := os.Stat(o.Filepath)
	if err != nil {
		return err
	}
	file, err := os.ReadFile(stat.Name())
	if err != nil {
		return err
	}

	template := ImportTemplate{}
	err = json.Unmarshal(file, &template)
	if err != nil {
		return err
	}

	o.create = &CreateOptions{
		WorkspaceName:    o.WorkspaceName,
		Name:             template.Name,
		Description:      template.Description,
		Language:         template.Language,
		Source:           template.Source,
		URL:              template.URL,
		Tag:              template.Tag,
		Token:            template.Token,
		MainWorkflowPath: template.MainWorkflowPath,
		options:          o.options,
	}
	if err = o.Complete(); err != nil {
		return err
	}
	if err = o.options.Validate(); err != nil {
		return err
	}
	if err = o.create.Complete(); err != nil {
		return err
	}
	return o.create.Validate()
}

type ImportTemplate struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Language         string `json:"language"`
	Source           string `json:"source"`
	URL              string `json:"URL"`
	Tag              string `json:"tag"`
	Token            string `json:"token"`
	MainWorkflowPath string `json:"mainWorkflowPath"`
}

// Run run the create workflow command
func (o *ImportOptions) Run(args []string) error {
	return o.create.Run([]string{o.create.Name})
}
