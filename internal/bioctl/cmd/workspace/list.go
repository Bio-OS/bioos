package workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/consts"
)

// ListOptions is an options to list workspaces.
type ListOptions struct {
	Page        int32
	Size        int32
	OrderBy     string
	SearchWords []string
	Ids         []string

	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewListOptions returns a reference to a ListOptions.
func NewListOptions(opt *clioptions.GlobalOptions) *ListOptions {
	return &ListOptions{
		options: opt,
	}
}

// NewCmdList new a list workspace cmd.
func NewCmdList(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewListOptions(opt)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list workspaces",
		Long:  "list workspaces",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().Int32VarP(&o.Page, "page", "p", 1, "The page number")
	cmd.Flags().Int32VarP(&o.Size, "size", "s", 10, "The page size")
	cmd.Flags().StringVar(&o.OrderBy, "order-by", o.OrderBy, "The order-by field")
	cmd.Flags().StringSliceVar(&o.SearchWords, "search-word", o.SearchWords, "The search word")
	cmd.Flags().StringSliceVar(&o.Ids, "ids", o.Ids, "The ids of the workspace.")

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
	return nil
}

// Run run the create workspace command
func (o *ListOptions) Run(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	req := &convert.ListWorkspacesRequest{
		Page: int(o.Page),
		Size: int(o.Size),
	}
	if len(o.Ids) > 0 {
		req.IDs = strings.Join(o.Ids, consts.QuerySliceDelimiter)
	}
	if o.OrderBy != "" {
		req.OrderBy = o.OrderBy
	}
	if len(o.SearchWords) > 0 {
		req.SearchWord = strings.Join(o.SearchWords, consts.QuerySliceDelimiter)
	}
	resp, err := o.workspaceClient.ListWorkspaces(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp)

	return nil
}

func (o *ListOptions) GetPromptArgs() ([]string, error) {
	return []string{}, nil
}

func (o *ListOptions) GetPromptOptions() error {
	var err error
	o.Page, err = prompt.PromptRequiredInt32("Page")
	if err != nil {
		return err
	}

	o.Size, err = prompt.PromptRequiredInt32("Size")
	if err != nil {
		return err
	}

	orderByFields, err := prompt.PromptStringMultiSelect("OrderBy", 2, []string{workspace.OrderByName, workspace.OrderByCreateTime})
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

	o.SearchWords, err = prompt.PromptStringSlice("SearchWords")
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

	return err
}

func (o *ListOptions) GetDefaultFormat() formatter.Format {
	return formatter.TableFormat
}

// ConvertWorkspaceNameIntoID convert workspace name into id
func ConvertWorkspaceNameIntoID(ctx context.Context, workspaceClient factory.WorkspaceClient, name string) (id string, err error) {
	pageNum := 1
	pageSize := 100
	for {
		var workspaces *convert.ListWorkspacesResponse
		workspaces, err = workspaceClient.ListWorkspaces(ctx, &convert.ListWorkspacesRequest{
			Page:       pageNum,
			Size:       pageSize,
			SearchWord: name,
		})
		if err != nil {
			return
		}
		for _, item := range workspaces.Items {
			if item.Name == name {
				id = item.Id
				return
			}
		}
		if len(workspaces.Items) < pageSize {
			err = fmt.Errorf("workspace: %s not found", name)
			return
		}
		pageNum++
	}
}
