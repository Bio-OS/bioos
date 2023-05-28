package data_model

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

// DeleteOptions is an options to delete a data-model.
type DeleteOptions struct {
	WorkspaceName string
	Name          string

	workspaceClient factory.WorkspaceClient
	dataModelClient factory.DataModelClient
	formatter       formatter.Formatter

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
		Use:   "delete",
		Short: "delete a data-model",
		Long:  "delete a data-model",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.Name, "name", "n", o.Name, "data-model names to List")

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

// Run run the delete data-model command
func (o *DeleteOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	req := &convert.ListDataModelsRequest{
		WorkspaceID: workspaceID,
	}

	remnant := make([]string, 0)
	if o.Name != "" {
		req.SearchWord = o.Name
		remnant = []string{o.Name}
	}

	listResp, err := o.dataModelClient.ListDataModels(ctx, req)
	if err != nil {
		return err
	}

	for _, item := range listResp.Items {
		// Delete all data-models
		if o.Name == "" {
			_, err := o.dataModelClient.DeleteDataModel(ctx, &convert.DeleteDataModelRequest{
				ID:          item.ID,
				WorkspaceID: workspaceID,
			})
			if err != nil {
				remnant = append(remnant, item.Name)
			}
			// Delete one specified data_model
		} else {
			if item.Name == o.Name {
				_, err := o.dataModelClient.DeleteDataModel(ctx, &convert.DeleteDataModelRequest{
					ID:          item.ID,
					WorkspaceID: workspaceID,
				})
				if err != nil {
					return err
				} //remove specified data_model name from remnant
				remnant = []string{}
				break
			}
		}
	}

	if len(remnant) > 0 {
		return fmt.Errorf("%v delete failed", remnant)
	}
	return nil

}

func (o *DeleteOptions) GetPromptArgs() ([]string, error) {
	return nil, nil
}

func (o *DeleteOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}

	o.Name, err = prompt.PromptOptionalString("DataModel Name", prompt.WithInputMessage("all dataModel will be deleted if no name specified"))
	if err != nil {
		return err
	}
	return nil
}

func (o *DeleteOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
