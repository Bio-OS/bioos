package workflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/pkg/consts"
)

// ListOptions is an options to list workflow.
type ListOptions struct {
	WorkspaceName string
	Page          int32
	Size          int32
	OrderBy       string
	SearchWords   []string
	Ids           []string

	workflowClient  factory.WorkflowClient
	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

func (o *ListOptions) GetPromptArgs() ([]string, error) {
	return []string{}, nil
}

func (o *ListOptions) GetPromptOptions() error {
	var err error
	o.WorkspaceName, err = prompt.PromptRequiredString("Workspace", prompt.WithInputMessage("Name"))
	if err != nil {
		return err
	}
	o.SearchWords, err = prompt.PromptStringSlice("SearchWords")
	if err != nil {
		return err
	}
	orderByFields, err := prompt.PromptStringMultiSelect("OrderBy", 2, []string{workflow.OrderByName, workflow.OrderByCreateTime})
	if err != nil {
		return err
	}
	if len(orderByFields) > 0 {
		for i, field := range orderByFields {
			ascending, err := prompt.PromptStringSelect(fmt.Sprintf("%s Ascending", field), 2, []string{consts.ASCOrdering, consts.DESCOrdering})
			if err != nil {
				return err
			}
			orderByFields[i] += consts.OrderDelimiter + ascending
		}
		o.OrderBy = strings.Join(orderByFields, ",")
	}

	o.Page, err = prompt.PromptRequiredInt32("Page")
	if err != nil {
		return err
	}
	o.Size, err = prompt.PromptRequiredInt32("Size")
	if err != nil {
		return err
	}
	o.Ids, err = prompt.PromptStringSlice("IDs")
	if err != nil {
		return err
	}
	// correct formatter
	if len(o.Ids) > 0 {
		o.options.Stream.OutputFormat = formatter.JsonFormat
		o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	}

	return nil
}

func (o *ListOptions) GetDefaultFormat() formatter.Format {
	return formatter.TableFormat
}

// NewListOptions returns a reference to a ListOptions.
func NewListOptions(opt *clioptions.GlobalOptions) *ListOptions {
	return &ListOptions{
		options: opt,
	}
}

// NewCmdList new a list workflow cmd.
func NewCmdList(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewListOptions(opt)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list workflow",
		Long:  "list workflow",
		Args:  cobra.ExactArgs(0),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The name of the workspace.")
	cmd.Flags().Int32Var(&o.Page, "page", 1, "The page number")
	cmd.Flags().Int32Var(&o.Size, "size", 10, "The page size")
	cmd.Flags().StringVar(&o.OrderBy, "order-by", o.OrderBy, "The order-by field: Name:desc")
	cmd.Flags().StringSliceVar(&o.SearchWords, "search-words", o.SearchWords, "The search word")
	cmd.Flags().StringSliceVar(&o.Ids, "ids", o.Ids, "The ids of the workspace.")

	return cmd
}

// Complete completes all the required options.
func (o *ListOptions) Complete() error {
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
		if len(o.Ids) > 0 {
			o.options.Stream.OutputFormat = formatter.JsonFormat
		} else {
			o.options.Stream.OutputFormat = o.GetDefaultFormat()
		}
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the list options
func (o *ListOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.WorkspaceName == "" {
		return fmt.Errorf("workspace not provide")
	}
	return nil
}

// Run run the create workflow command
func (o *ListOptions) Run(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}
	req := &convert.ListWorkflowsRequest{
		Page:        int(o.Page),
		Size:        int(o.Size),
		WorkspaceID: workspaceID,
	}
	if o.OrderBy != "" {
		req.OrderBy = o.OrderBy
	}
	if len(o.SearchWords) > 0 {
		req.SearchWord = strings.Join(o.SearchWords, consts.QuerySliceDelimiter)
	}
	if len(o.Ids) > 0 {
		req.IDs = strings.Join(o.Ids, consts.QuerySliceDelimiter)
	}
	resp, err := o.workflowClient.ListWorkflow(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp)

	return nil
}

// ConvertWorkflowNameIntoID convert workflow name into id
func ConvertWorkflowNameIntoID(ctx context.Context, workflowClient factory.WorkflowClient, wsID string, name string) (id string, err error) {
	pageNum := 1
	pageSize := 100
	for {
		var listWorkflowsResponse *convert.ListWorkflowsResponse
		listWorkflowsResponse, err = workflowClient.ListWorkflow(ctx, &convert.ListWorkflowsRequest{
			Page:        pageNum,
			Size:        pageSize,
			SearchWord:  name,
			WorkspaceID: wsID,
		})
		if err != nil {
			return
		}
		for _, item := range listWorkflowsResponse.Items {
			if item.Name == name {
				id = item.ID
				return
			}
		}
		if len(listWorkflowsResponse.Items) < pageSize {
			err = fmt.Errorf("workflow: %s not found", name)
			return
		}
		pageNum++
	}
}
