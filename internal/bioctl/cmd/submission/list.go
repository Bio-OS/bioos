package submission

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/cmd"
	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/pkg/consts"
)

// ListOptions is an options to List submissions.
type ListOptions struct {
	WorkspaceName string
	Status        []string
	Page          int32
	Size          int32
	OrderBy       string
	SearchWords   []string
	Ids           []string

	workspaceClient  factory.WorkspaceClient
	submissionClient factory.SubmissionClient
	formatter        formatter.Formatter

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
		Short: "list submissions",
		Long:  "list submissions of a specified workspace",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringSliceVar(&o.Status, "status", o.Status, "filtered status")
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
	o.submissionClient, err = f.SubmissionClient()
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

// Validate validate the List options
func (o *ListOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}
	return nil
}

// Run run the List submission command
func (o *ListOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	req := &convert.ListSubmissionsRequest{
		WorkspaceID: workspaceID,
		Page:        int(o.Page),
		Size:        int(o.Size),
	}

	if o.OrderBy != "" {
		req.OrderBy = o.OrderBy
	}
	if len(o.SearchWords) > 0 {
		req.SearchWord = strings.Join(o.SearchWords, consts.QuerySliceDelimiter)
	}
	if len(o.Status) > 0 {
		req.Status = o.Status
	}
	if len(o.Ids) > 0 {
		req.IDs = o.Ids
	}
	resp, err := o.submissionClient.ListSubmissions(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp)

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
	o.Page, err = prompt.PromptRequiredInt32("Page")
	if err != nil {
		return err
	}

	o.Size, err = prompt.PromptRequiredInt32("Size")
	if err != nil {
		return err
	}

	orderByFields, err := prompt.PromptStringMultiSelect("OrderBy", 2, []string{submission.OrderByName, submission.OrderByStartTime})
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

	o.Status, err = prompt.PromptStringMultiSelect("Status", 7, []string{consts.SubmissionPending,
		consts.SubmissionRunning, consts.SubmissionFailed, consts.SubmissionFinished, consts.SubmissionSucceeded, consts.SubmissionCancelling,
		consts.SubmissionCancelled})
	if err != nil {
		return err
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

	return nil
}

func (o *ListOptions) GetDefaultFormat() formatter.Format {
	return formatter.TableFormat
}
