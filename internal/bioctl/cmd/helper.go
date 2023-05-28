package cmd

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

func GetAllRuns(ctx context.Context, submissionClient factory.SubmissionClient, workspaceID string, submissionID string) ([]convert.RunItem, error) {
	resp, err := submissionClient.ListRuns(ctx, &convert.ListRunsRequest{
		Page:         1,
		Size:         100,
		WorkspaceID:  workspaceID,
		SubmissionID: submissionID,
	})
	if err != nil {
		return nil, err
	}
	currentNum := resp.Size
	currentPage := resp.Page

	items := resp.Items

	for currentNum < resp.Total {
		currentNum += len(resp.Items)
		currentPage += 1

		resp, err = submissionClient.ListRuns(ctx, &convert.ListRunsRequest{
			Page:         currentPage,
			Size:         100,
			WorkspaceID:  workspaceID,
			SubmissionID: submissionID,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, resp.Items...)
	}
	return items, nil
}

func ParseWholeDataModel(ctx context.Context, dataModelClient factory.DataModelClient, workspaceID, dataModelID string) ([]string, [][]string, error) {
	resp, err := dataModelClient.ListDataModelRows(ctx, &convert.ListDataModelRowsRequest{
		Page:        1,
		Size:        100,
		WorkspaceID: workspaceID,
		ID:          dataModelID,
	})
	if err != nil {
		return nil, nil, err
	}
	currentNum := int64(resp.Size)
	currentPage := resp.Page

	headers := resp.Headers
	rows := resp.Rows

	for currentNum < resp.Total {
		currentNum += int64(resp.Size)
		currentPage += 1

		resp, err := dataModelClient.ListDataModelRows(ctx, &convert.ListDataModelRowsRequest{
			Page:        currentPage,
			Size:        100,
			WorkspaceID: workspaceID,
			ID:          dataModelID,
		})
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, resp.Rows...)
	}
	return headers, rows, nil
}

func GetWorkspaceName(timeout int, workspaceClient factory.WorkspaceClient) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	resp, err := workspaceClient.ListWorkspaces(ctx, &convert.ListWorkspacesRequest{})
	if err != nil {
		return "", err
	}

	if resp.Total == 0 {
		return "", fmt.Errorf("no workspace found")
	}

	var workspaceName string
	if resp.Total > 20 {
		workspaceName, err = prompt.PromptRequiredString("WorkspaceName", prompt.WithInputMessage(""))
	} else {
		items := make([]string, len(resp.Items))
		for j := range resp.Items {
			workspace := resp.Items[j]
			items[j] = fmt.Sprintf("%s", workspace.Name)
		}
		workspaceName, err = prompt.PromptStringSelect("WorkspaceName", 10, items)
	}
	if err != nil {
		return "", err
	}
	return workspaceName, nil
}

var ExampleTemplate = template.Must(template.New("Example").Parse(
	`	# Use Command-Line mode:
	"bioctl {{.Use}} {{- range .OptionsUsage }} {{ . }} {{- end }}"
	
	# Enter Interactive mode:
	1. Firstly input parent Command:
	"bioctl {{.ParentName}}" 

	2. Select Subcommand 
	"{{.Name}}"
`))

type ExampleData struct {
	ParentName   string
	Name         string
	Use          string
	OptionsUsage []string
}

// "# list workspaces: bioctl workspace list"
func GenExample(cmd *cobra.Command) string {
	// template.
	flags := cmd.Flags()
	// cmd.Args
	flags.FlagUsages()
	data := ExampleData{
		ParentName: cmd.Parent().Name(),
		Name:       cmd.Name(),
		Use:        cmd.Use,
	}
	flags.VisitAll(func(flag *pflag.Flag) {
		var option string
		if flag.Shorthand != "" {
			option += fmt.Sprintf("-%s/", flag.Shorthand)
		}
		option += fmt.Sprintf("--%s [%s]", flag.Name, flag.Value.Type())
		data.OptionsUsage = append(data.OptionsUsage, option)
	})

	var res bytes.Buffer
	ExampleTemplate.Execute(&res, data)

	return res.String()
}
