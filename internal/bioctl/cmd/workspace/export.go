package workspace

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/schema"
	"github.com/Bio-OS/bioos/pkg/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type ExportOptions struct {
	OutputPath  string
	WorkspaceID string

	workspaceClient factory.WorkspaceClient
	dataModelClient factory.DataModelClient
	workflowClient  factory.WorkflowClient
	notebookClient  factory.NotebookClient

	formatter formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewExportOptions returns a reference to a ExportOptions
func NewExportOptions(opt *clioptions.GlobalOptions) *ExportOptions {
	return &ExportOptions{
		options: opt,
	}
}

func NewCmdExport(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewExportOptions(opt)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "export a workspace",
		Long:  "export a workspace",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.OutputPath, "output", "p", o.OutputPath, "The output path of the workspace zip file")
	cmd.Flags().StringVarP(&o.WorkspaceID, "workspaceID", "w", o.WorkspaceID, "The id of the workspace.")

	return cmd
}

// Complete completes all the required options.
func (o *ExportOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workspaceClient, err = f.WorkspaceClient()
	if err != nil {
		return err
	}
	o.workflowClient, err = f.WorkflowClient()
	if err != nil {
		return err
	}
	o.notebookClient, err = f.NotebookClient()
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

// Validate validate the export options
func (o *ExportOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.OutputPath == "" {
		return fmt.Errorf("you must choose a output path to export workspace and zip")
	}
	_, err := os.Stat(o.OutputPath)
	if err != nil {
		return err
	}
	if len(o.WorkspaceID) == 0 {
		return fmt.Errorf("workspace id [%s] should not be empty", o.WorkspaceID)
	}
	return nil
}

// Run run the import workspace command
func (o *ExportOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	// set umask 0, set the default permissions given to a file
	mask := syscall.Umask(0)
	// recover umask
	defer syscall.Umask(mask)

	workspaceTypedSchema, workspaceDir, err := o.exportWorkspaceToSchema(ctx, o.OutputPath)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(workspaceTypedSchema)
	if err != nil {
		return fmt.Errorf("failed to marshal workspace typed schema err: %w", err)
	}
	workspaceFilePath := path.Join(workspaceDir, consts.WorkspaceYAMLName)
	err = os.WriteFile(workspaceFilePath, bytes, consts.SchemaFileMode)
	if err != nil {
		return fmt.Errorf("failed to save workflow yaml file err: %w", err)
	}

	o.formatter.Write(workspaceDir)

	return nil
}

func (o *ExportOptions) GetPromptArgs() ([]string, error) {
	return nil, nil
}

func (o *ExportOptions) GetPromptOptions() error {
	var err error

	o.OutputPath, err = prompt.PromptStringWithValidator("Workspace Output Path", func(ans interface{}) error {
		err = survey.Required(ans)
		if err != nil {
			return err
		}
		_, err = os.Stat(cast.ToString(ans))
		if err != nil {
			curPath, _ := os.Getwd()
			return fmt.Errorf("%w, (currenct path is [%s])", err, curPath)
		}
		return nil
	}, prompt.WithInputMessage("Workspace Exported Zip File Output Path"))

	if err != nil {
		return err
	}

	o.WorkspaceID, err = prompt.PromptRequiredString("Workspace ID", prompt.WithInputMessage("workspace ID"))
	if err != nil {
		return err
	}

	return nil
}

func (o *ExportOptions) exportWorkspaceToSchema(ctx context.Context, baseDir string) (*schema.WorkspaceTypedSchema, string, error) {
	workspace, workspaceTypedSchema, err := o.getSchemaFromWorkspace(ctx)
	if err != nil {
		return nil, "", err
	}
	workspaceDir, err := mkdirWorkspaceDir(baseDir, workspace.Name)
	if err != nil {
		return nil, "", err
	}

	type releaseFn func(context.Context, *convert.WorkspaceItem, *schema.WorkspaceTypedSchema, string) error
	fn := []releaseFn{o.exportNotebookToSchema, o.exportDataModelToSchema, o.exportWorkflowToSchema}
	for _, f := range fn {
		if err = f(ctx, workspace, workspaceTypedSchema, workspaceDir); err != nil {
			return nil, "", fmt.Errorf("export workspace schema fail: %w", err)
		}
	}
	return workspaceTypedSchema, workspaceDir, nil
}

func (o *ExportOptions) getSchemaFromWorkspace(ctx context.Context) (*convert.WorkspaceItem, *schema.WorkspaceTypedSchema, error) {
	workspace, err := o.workspaceClient.GetWorkspace(ctx, &convert.GetWorkspaceRequest{Id: o.WorkspaceID})
	if err != nil {
		return nil, nil, fmt.Errorf("get workspace from db failed: %w", err)
	}
	return &workspace.WorkspaceItem, &schema.WorkspaceTypedSchema{
		Name:        workspace.Name,
		Version:     consts.WorkspaceScopedSchemaVersion,
		Description: workspace.Description,
	}, nil
}

// jupyterhub data
// - save dashboard ipynb file
// - save notebook ipynb file
// - get image info from db
func (o *ExportOptions) exportNotebookToSchema(ctx context.Context, workspace *convert.WorkspaceItem, workspaceTypedSchema *schema.WorkspaceTypedSchema, workspaceDir string) error {
	notebookDir := path.Join(workspaceDir, consts.NotebookDirName)
	err := os.MkdirAll(notebookDir, consts.SchemaFileMode)
	if err != nil {
		return fmt.Errorf("failed create notebook dir: %w", err)
	}
	artifacts := make([]*schema.Artifact, 0)
	resp, err := o.notebookClient.ListNotebooks(ctx, &convert.ListNotebooksRequest{
		WorkspaceID: workspace.Id,
	})
	if err != nil {
		return fmt.Errorf("list notebooks error: %w", err)
	}
	for _, notebookItem := range resp.Items {
		notebookResp, err := o.notebookClient.GetNotebook(ctx, &convert.GetNotebookRequest{
			WorkspaceID: workspace.Id,
			Name:        notebookItem.Name,
		})
		if err != nil {
			return fmt.Errorf("get notebook error: %w", err)
		}
		filePath := path.Join(notebookDir, fmt.Sprintf("%s%s", notebookItem.Name, notebook.NotebookFileExt))
		err = writeFile(filePath, string(notebookResp.Content))
		if err != nil {
			return fmt.Errorf("write notebook's file[%s] error: %w", notebookItem.Name, err)
		}
		artifact := &schema.Artifact{
			Name: notebookItem.Name,
			Path: path.Join(consts.NotebookDirName, fmt.Sprintf("%s%s", notebookItem.Name, notebook.NotebookFileExt)),
		}
		artifacts = append(artifacts, artifact)
	}

	workspaceTypedSchema.Notebooks = schema.NotebookTypedSchema{
		Artifacts: artifacts,
	}
	return nil
}

// dataModel data
// - get data model from db
// - get data model rows from db
// - export .csv file
func (o *ExportOptions) exportDataModelToSchema(ctx context.Context, workspace *convert.WorkspaceItem, workspaceTypedSchema *schema.WorkspaceTypedSchema, workspaceDir string) error {
	dataModelDir := path.Join(workspaceDir, consts.DataModelDirName)
	err := os.MkdirAll(dataModelDir, consts.SchemaFileMode)
	if err != nil {
		return fmt.Errorf("failed create data model dir: %w", err)
	}
	dataModelTypedSchemas := make([]schema.DataModelTypedSchema, 0)
	resp, err := o.dataModelClient.ListDataModels(ctx, &convert.ListDataModelsRequest{
		WorkspaceID: workspace.Id,
	})
	if err != nil {
		return fmt.Errorf("list source data models error: %w", err)
	}
	for _, dataModel := range resp.Items {
		headers, rows, err := o.getAllDataModelRows(ctx, workspace, dataModel)
		if err != nil {
			return fmt.Errorf("list source data model's rows error: %w", err)
		}
		if err = utils.WriteDataModelToCSVFile(path.Join(dataModelDir, fmt.Sprintf("%s.csv", dataModel.Name)), headers, rows); err != nil {
			return fmt.Errorf("export data model csv file error: %w", err)
		}
		dataModelTypedSchema := schema.DataModelTypedSchema{
			Name: dataModel.Name,
			Type: getDataModelType(utils.GetDataModelType(dataModel.Name)),
			Path: path.Join(consts.DataModelDirName, fmt.Sprintf("%s.csv", dataModel.Name)),
		}
		dataModelTypedSchemas = append(dataModelTypedSchemas, dataModelTypedSchema)
	}
	workspaceTypedSchema.DataModels = dataModelTypedSchemas
	return nil
}

// workflow data
// - get from db
// - save wdl file from s3
func (o *ExportOptions) exportWorkflowToSchema(ctx context.Context, workspace *convert.WorkspaceItem, workspaceTypedSchema *schema.WorkspaceTypedSchema, workspaceDir string) error {
	workflowDir := path.Join(workspaceDir, consts.WorkflowDirName)
	err := os.MkdirAll(workflowDir, consts.SchemaFileMode)
	if err != nil {
		return fmt.Errorf("failed create workflow dir: %w", err)
	}
	workflows, err := o.listAllWorkflows(ctx, workspace)
	if err != nil {
		return fmt.Errorf("list all workflows error: %w", err)
	}
	workflowTypedSchemas := make([]schema.WorkflowTypedSchema, 0)
	if err != nil {
		return fmt.Errorf("failed to list source workflows: %w", err)
	}
	for _, workflow := range workflows {
		if workflow.LatestVersion.Status != "Success" {
			continue
		}
		fileResp, err := o.workflowClient.ListWorkflowFiles(ctx, &convert.ListWorkflowFilesRequest{
			WorkflowID:        workflow.ID,
			WorkspaceID:       workspace.Id,
			WorkflowVersionID: workflow.LatestVersion.ID,
			Page:              1,
			Size:              100,
		})
		if err != nil {
			return fmt.Errorf("list workflow's files error: %w", err)
		}
		for _, file := range fileResp.Items {
			decodedContent, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				return fmt.Errorf("base64 decode error: %w", err)
			}
			err = writeFile(path.Join(workflowDir, workflow.Name, file.Path), string(decodedContent))
			if err != nil {
				return fmt.Errorf("write workflow's file[%s] error: %w", file.Path, err)
			}
		}
		workflowTypedSchema := schema.WorkflowTypedSchema{
			Name:             workflow.Name,
			Description:      utils.PointString(workflow.Description),
			Language:         workflow.LatestVersion.Language,
			Version:          utils.PointString(workflow.LatestVersion.LanguageVersion),
			MainWorkflowPath: workflow.LatestVersion.MainWorkflowPath,
			Path:             path.Join(consts.WorkflowDirName, workflow.Name),
		}
		if workflow.LatestVersion.Source == "git" {
			workflowTypedSchema.Metadata = schema.WorkflowMetadata{
				Tag:   workflow.LatestVersion.Metadata["gitTag"],
				Token: utils.PointString(workflow.LatestVersion.Metadata["gitToken"]),
			}
			gitURL, ok := workflow.LatestVersion.Metadata["gitURL"]
			if ok {
				parsedURL, err := url.Parse(gitURL)
				if err != nil {
					return fmt.Errorf("parse workflow's git url error: %w", err)
				}
				workflowTypedSchema.Metadata.Scheme = parsedURL.Scheme
				workflowTypedSchema.Metadata.Repo = strings.TrimPrefix(gitURL, fmt.Sprintf("%s://", parsedURL.Scheme))
			}
		}
		workflowTypedSchemas = append(workflowTypedSchemas, workflowTypedSchema)
	}
	workspaceTypedSchema.Workflows = workflowTypedSchemas
	return nil
}

func (o *ExportOptions) getAllDataModelRows(ctx context.Context, workspace *convert.WorkspaceItem, dataModel convert.DataModel) (headers []string, rows [][]string, err error) {
	page := int32(1)
	size := int32(100)
	req := &convert.ListDataModelRowsRequest{
		WorkspaceID: workspace.Id,
		ID:          dataModel.ID,
		Page:        page,
		Size:        size,
	}
	resp, err := o.dataModelClient.ListDataModelRows(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	headers = resp.Headers
	rows = resp.Rows
	total := resp.Total
	if total <= int64(resp.Size) {
		return headers, rows, nil
	}
	breakFlag := false
	for {
		page++
		if total <= int64(page*size) {
			size = int32(total - int64((page-1)*size))
			breakFlag = true
		}
		req.Page = page
		req.Size = size
		resp, err = o.dataModelClient.ListDataModelRows(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, resp.Rows...)
		if breakFlag {
			break
		}
	}
	return headers, rows, nil
}

func (o *ExportOptions) listAllWorkflows(ctx context.Context, workspace *convert.WorkspaceItem) ([]convert.WorkflowItem, error) {
	page := 1
	size := 100
	req := &convert.ListWorkflowsRequest{
		WorkspaceID: workspace.Id,
		OrderBy:     "name:asc",
		Page:        page,
		Size:        size,
	}
	resp, err := o.workflowClient.ListWorkflow(ctx, req)
	if err != nil {
		return nil, err
	}
	workflows := resp.Items
	total := resp.Total
	if total <= resp.Size {
		return workflows, nil
	}
	breakFlag := false
	for {
		page++
		if total <= page*size {
			size = total - (page-1)*size
			breakFlag = true
		}
		req.Page = page
		req.Size = size
		resp, err = o.workflowClient.ListWorkflow(ctx, req)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, resp.Items...)
		if breakFlag {
			break
		}
	}
	return workflows, nil
}

func getDataModelType(dbType string) string {
	switch dbType {
	case consts.DataModelTypeEntitySet:
		return "entitySet"
	case consts.DataModelTypeEntity:
		return "entity"
	case consts.DataModelTypeWorkspace:
		return "workspace"
	default:
		return "entity"
	}
}

func writeFile(filePath, content string) error {
	fileDir := path.Dir(filePath)
	if err := os.MkdirAll(fileDir, consts.SchemaFileMode); err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() // nolint

	if _, err = file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func mkdirWorkspaceDir(baseDir, workspaceName string) (workspaceDir string, err error) {
	var overwriteCommand string
	workspaceDir = path.Join(baseDir, workspaceName)
	if _, err = os.Stat(workspaceDir); err == nil {
		fmt.Printf("allow overwrite on %s? yes or no\n", workspaceDir)
		if _, err = fmt.Scanln(&overwriteCommand); err != nil {
			return workspaceDir, fmt.Errorf("failed to scan:%w", err)
		} else if overwriteCommand == "yes" {
			err := os.RemoveAll(workspaceDir)
			if err != nil {
				return workspaceDir, fmt.Errorf("failed to delete dir err: %w", err)
			}
		} else {
			return workspaceDir, fmt.Errorf("%s exist and do not overwrite", workspaceDir)
		}
	} else if !os.IsNotExist(err) {
		return workspaceDir, fmt.Errorf("failed to stat file [%s]:%w", workspaceDir, err)
	}

	err = os.MkdirAll(workspaceDir, consts.SchemaFileMode)
	if err != nil {
		return workspaceDir, fmt.Errorf("failed to create workspace dir err: %w", err)
	}
	return workspaceDir, nil
}

func (o *ExportOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
