package data_model

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/cmd"
	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// ImportOptions is an options to import a data-model.
type ImportOptions struct {
	WorkspaceName string
	InputFile     string

	workspaceClient factory.WorkspaceClient
	dataModelClient factory.DataModelClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewImportOptions returns a reference to a ImportOptions
func NewImportOptions(opt *clioptions.GlobalOptions) *ImportOptions {
	return &ImportOptions{
		options: opt,
	}
}

func NewCmdImport(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewImportOptions(opt)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "import a data-model",
		Long:  "import a data-model",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.InputFile, "input-file ", "i", o.InputFile, "the file (only support csv) to import")

	return cmd
}

// Complete completes all the required options.
func (o *ImportOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workspaceClient, err = f.WorkspaceClient()
	if err != nil {
		return err
	}
	o.dataModelClient, err = f.DataModelClient()
	if err != nil {
		return err
	}

	if o.options.Stream.OutputFormat == "" {
		o.options.Stream.OutputFormat = o.GetDefaultFormat()
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the import options
func (o *ImportOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}
	if o.InputFile == "" || path.Ext(o.InputFile) != ".csv" {
		return fmt.Errorf("you must choose a csv file to import")
	}
	_, err := os.Stat(o.InputFile)
	if err != nil {
		return err
	}
	return nil
}

// Run run the import data-model command
func (o *ImportOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}
	name := strings.TrimSuffix(path.Base(o.InputFile), path.Ext(o.InputFile))

	listResp, err := o.dataModelClient.ListDataModels(ctx, &convert.ListDataModelsRequest{
		WorkspaceID: workspaceID,
		Types:       []string{consts.DataModelTypeEntity},
		SearchWord:  name,
	})
	if err != nil {
		return err
	}
	for _, item := range listResp.Items {
		if item.Name == name {
			_, err = o.dataModelClient.DeleteDataModel(ctx, &convert.DeleteDataModelRequest{
				ID:          item.ID,
				WorkspaceID: workspaceID,
			})
			if err != nil {
				return fmt.Errorf("delete existed dataModel %s failed: %w", name, err)
			}
		}
	}

	req := &convert.PatchDataModelRequest{
		WorkspaceID: workspaceID,
		Name:        name,
	}

	headers, rows, err := utils.ReadDataModelFromCSV(o.InputFile)
	if err != nil {
		return err
	}
	req.Headers = headers
	req.Rows = make([][]string, len(rows))
	for i, data := range rows {
		req.Rows[i] = data
	}

	resp, err := o.dataModelClient.PatchDataModel(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp.ID)

	return nil
}

func (o *ImportOptions) GetPromptArgs() ([]string, error) {
	return nil, nil
}

func (o *ImportOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}

	o.InputFile, err = prompt.PromptStringWithValidator("Input File Path", func(ans interface{}) error {
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
	}, prompt.WithInputMessage("only support CSV"))

	if err != nil {
		return err
	}

	return nil
}

func (o *ImportOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
