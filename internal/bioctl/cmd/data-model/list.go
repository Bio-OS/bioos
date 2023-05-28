package data_model

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/cmd"
	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/pkg/consts"
)

// ListOptions is an options to List data-models.
type ListOptions struct {
	WorkspaceName string
	Types         []string
	Name          string

	workspaceClient factory.WorkspaceClient
	dataModelClient factory.DataModelClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewListOptions returns a reference to a ListOptions
func NewListOptions(opt *clioptions.GlobalOptions) *ListOptions {
	return &ListOptions{
		options: opt,
	}
}

func NewCmdList(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewListOptions(opt)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list data-models",
		Long:  "list data-models of a specified workspace",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.Name, "name", "n", o.Name, "data-model names to List")
	cmd.Flags().StringSliceVarP(&o.Types, "types", "t", o.Types, "data-model types to List")

	return cmd
}

// Complete completes all the required options.
func (o *ListOptions) Complete() error {
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

// Validate validate the List options
func (o *ListOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}
	for _, t := range o.Types {
		if t != consts.DataModelTypeEntity && t != consts.DataModelTypeEntitySet && t != consts.DataModelTypeWorkspace {
			return fmt.Errorf("data-model type %s not support", t)
		}
	}
	return nil
}

// Run run the List data-model command
func (o *ListOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	if o.Name != "" {
		dataModelID, err := ConvertDataModelNameIntoID(ctx, o.dataModelClient, workspaceID, o.Name)
		if err != nil {
			return err
		}
		headers, rows, err := cmd.ParseWholeDataModel(ctx, o.dataModelClient, workspaceID, dataModelID)
		if err != nil {
			return err
		}
		if len(rows) > 10000 {
			o.formatter.Write(fmt.Sprintf("It will take a while to collect in that the length of dataModel is too big: %d rows", len(rows)))
		}
		o.formatter.Write(getDataModelStruct(headers, rows))
		return nil
	}

	req := &convert.ListDataModelsRequest{
		WorkspaceID: workspaceID,
	}
	if len(o.Types) > 0 {
		req.Types = o.Types
	}
	resp, err := o.dataModelClient.ListDataModels(ctx, req)
	if err != nil {
		return err
	}

	o.formatter.Write(resp.Items)
	return nil
}

func (o *ListOptions) GetPromptArgs() ([]string, error) {
	return nil, nil
}

func (o *ListOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}
	o.Name, err = prompt.PromptOptionalString("DataModel Name")
	if err != nil {
		return err
	}
	if o.Name == "" {
		o.Types, err = prompt.PromptStringMultiSelect("DataModel Types", 3,
			[]string{consts.DataModelTypeEntity, consts.DataModelTypeEntitySet, consts.DataModelTypeWorkspace},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ListOptions) GetDefaultFormat() formatter.Format {
	return formatter.TableFormat
}

func ConvertDataModelNameIntoID(ctx context.Context, dataModelClient factory.DataModelClient, wsID, dataModelName string) (string, error) {
	resp, err := dataModelClient.ListDataModels(ctx, &convert.ListDataModelsRequest{
		WorkspaceID: wsID,
		SearchWord:  dataModelName,
	})
	if err != nil {
		return "", err
	}

	for _, item := range resp.Items {
		if item.Name == dataModelName {
			return item.ID, nil
		}
	}
	return "", fmt.Errorf("no data-model named %s found", dataModelName)
}

func getDataModelStruct(headers []string, rows [][]string) interface{} {
	f := make([]reflect.StructField, len(headers))
	fields := headers
	tOfStr := reflect.TypeOf("")

	for i, v := range fields {
		f[i] = reflect.StructField{
			Name: strcase.ToCamel(reflect.ValueOf(v).Interface().(string)),
			Type: tOfStr,
		}
	}

	t := reflect.StructOf(f)
	items := reflect.MakeSlice(reflect.SliceOf(t), len(rows), len(rows))

	for i, row := range rows {
		e := reflect.New(t).Elem()
		for j, elem := range row {
			e.Field(j).Set(reflect.ValueOf(elem))
		}
		items.Index(i).Set(e)
	}

	return items.Interface()
}
