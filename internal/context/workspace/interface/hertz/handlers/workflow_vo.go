package handlers

import (
	"strings"
	"time"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workflow"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type createWorkflowRequest struct {
	ID               string  `json:"id"`
	WorkspaceID      string  `path:"workspace-id"`
	Name             string  `json:"name" validate:"required,resName"`
	Description      *string `json:"description" validate:"workspaceDesc"`
	Language         string  `json:"language" validate:"required,oneof=WDL NextFlow"`
	Source           string  `json:"source" validate:"required,oneof=git"`
	URL              string  `json:"url" validate:"required"`
	Tag              string  `json:"tag" validate:"required"`
	Token            string  `json:"token"`
	MainWorkflowPath string  `json:"mainWorkflowPath" validate:"required"`
}

func (req *createWorkflowRequest) toDTO() *command.CreateCommand {
	return &command.CreateCommand{
		ID:               req.ID,
		Name:             req.Name,
		Description:      req.Description,
		WorkspaceID:      req.WorkspaceID,
		Language:         req.Language,
		Source:           req.Source,
		URL:              req.URL,
		Tag:              req.Tag,
		Token:            req.Token,
		MainWorkflowPath: req.MainWorkflowPath,
	}
}

type createWorkflowResponse struct {
	ID string `json:"id"`
}

type getWorkflowRequest struct {
	WorkspaceID string `path:"workspace-id"`
	ID          string `path:"id"`
}

func (req *getWorkflowRequest) toDTO() *query.GetQuery {
	return &query.GetQuery{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
	}
}

type getWorkflowResponse struct {
	Workflow *WorkflowItem `json:"workflow"`
}

type listWorkflowsRequest struct {
	Page        int    `query:"page"`
	Size        int    `query:"size"`
	OrderBy     string `query:"orderBy"`
	SearchWord  string `query:"searchWord"`
	Exact       bool   `query:"exact"`
	IDs         string `query:"ids"`
	WorkspaceID string `path:"workspace-id"`
}

func (req listWorkflowsRequest) toDTO() (*query.ListQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}

	filter := &query.ListWorkflowsFilter{}
	if req.SearchWord != "" {
		filter.SearchWord = req.SearchWord
		filter.Exact = req.Exact
	}
	if len(req.IDs) > 0 {
		filter.IDs = strings.Split(req.IDs, consts.QuerySliceDelimiter)
	}
	return &query.ListQuery{
		Pg:          pg,
		Filter:      filter,
		WorkspaceID: req.WorkspaceID,
	}, nil
}

type listWorkflowsResponse struct {
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	Items []*WorkflowItem `json:"items"`
}

type updateWorkflowRequest struct {
	WorkspaceID      string  `path:"workspace-id"`
	ID               string  `path:"id"`
	Name             *string `json:"name,omitempty"`
	Description      *string `json:"description" validate:"workspaceDesc"`
	Language         *string `json:"language,omitempty" validate:"required,oneof=WDL"`
	Source           *string `json:"source,omitempty" validate:"required,oneof=git"`
	URL              *string `json:"url,omitempty"`
	Tag              *string `json:"tag,omitempty"`
	Token            *string `json:"token,omitempty"`
	MainWorkflowPath *string `json:"mainWorkflowPath,omitempty"`
}

func (req updateWorkflowRequest) toDTO() *command.UpdateCommand {
	reqDTO := &command.UpdateCommand{
		WorkspaceID:      req.WorkspaceID,
		ID:               req.ID,
		Name:             req.Name,
		Description:      req.Description,
		Language:         req.Language,
		Source:           req.Source,
		URL:              req.URL,
		Tag:              req.Tag,
		Token:            req.Token,
		MainWorkflowPath: req.MainWorkflowPath,
	}

	return reqDTO
}

type updateWorkflowResponse struct {
}

type deleteWorkflowRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
}

func (req deleteWorkflowRequest) toDTO() *command.DeleteCommand {
	return &command.DeleteCommand{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
	}
}

type deleteWorkflowResponse struct {
}

type getWorkflowFileRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
	WorkflowID  string `path:"workflow-id"`
}

func (req getWorkflowFileRequest) toDTO() *query.GetFileQuery {
	return &query.GetFileQuery{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
	}
}

type getWorkflowFileResponse struct {
	File *WorkflowFile `json:"file"`
}

type getWorkflowVersionRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
	WorkflowID  string `path:"workflow-id"`
}

func (req getWorkflowVersionRequest) toDTO() *query.GetVersionQuery {
	return &query.GetVersionQuery{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
	}
}

type getWorkflowVersionResponse struct {
	Version *WorkflowVersion `json:"version"`
}

type listWorkflowFilesRequest struct {
	Page              int    `query:"page"`
	Size              int    `query:"size"`
	OrderBy           string `query:"orderBy"`
	IDs               string `query:"ids"`
	WorkspaceID       string `path:"workspace-id"`
	WorkflowID        string `path:"workflow-id"`
	WorkflowVersionID string `query:"workflowVersionID,omitempty"`
}

func (req listWorkflowFilesRequest) toDTO() (*query.ListFilesQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}

	filter := &query.ListWorkflowFilesFilter{}

	if len(req.IDs) > 0 {
		filter.IDs = strings.Split(req.IDs, consts.QuerySliceDelimiter)
	}
	return &query.ListFilesQuery{
		Pg:                pg,
		Filter:            filter,
		WorkspaceID:       req.WorkspaceID,
		WorkflowID:        req.WorkflowID,
		WorkflowVersionID: req.WorkflowVersionID,
	}, nil
}

type listWorkflowFilesResponse struct {
	Page        int             `json:"page"`
	Size        int             `json:"size"`
	Total       int             `json:"total"`
	WorkspaceID string          `json:"workspaceID"`
	WorkflowID  string          `json:"workflowID"`
	Items       []*WorkflowFile `json:"items"`
}

type listWorkflowVersionsRequest struct {
	Page        int    `query:"page"`
	Size        int    `query:"size"`
	OrderBy     string `query:"orderBy"`
	IDs         string `query:"ids"`
	WorkspaceID string `path:"workspace-id"`
	WorkflowID  string `path:"workflow-id"`
}

func (req listWorkflowVersionsRequest) toDTO() (*query.ListVersionsQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}

	filter := &query.ListWorkflowVersionsFilter{}

	if len(req.IDs) > 0 {
		filter.IDs = strings.Split(req.IDs, consts.QuerySliceDelimiter)
	}
	return &query.ListVersionsQuery{
		Pg:          pg,
		Filter:      filter,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
	}, nil
}

type listWorkflowVersionsResponse struct {
	Page        int                `json:"page"`
	Size        int                `json:"size"`
	Total       int                `json:"total"`
	WorkspaceID string             `json:"workspaceID"`
	WorkflowID  string             `json:"workflowID"`
	Items       []*WorkflowVersion `json:"items"`
}

type WorkflowItem struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	LatestVersion *WorkflowVersion `json:"latestVersion"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

func WorkflowDTOtoWorkflowItemVO(workflowDTO *query.Workflow) *WorkflowItem {
	if workflowDTO == nil {
		return nil
	}
	return &WorkflowItem{
		ID:            workflowDTO.ID,
		Name:          workflowDTO.Name,
		Description:   workflowDTO.Description,
		LatestVersion: WorkflowVersionDTOtoVO(workflowDTO.LatestVersion),
		CreatedAt:     workflowDTO.CreatedAt,
		UpdatedAt:     workflowDTO.UpdatedAt,
	}
}

type WorkflowVersion struct {
	ID               string              `json:"id"`
	Status           string              `json:"status"`
	Message          string              `json:"message"`
	Language         string              `json:"language"`
	LanguageVersion  string              `json:"languageVersion"`
	MainWorkflowPath string              `json:"mainWorkflowPath"`
	Inputs           []*WorkflowParam    `json:"inputs"`
	Outputs          []*WorkflowParam    `json:"outputs"`
	Graph            string              `json:"graph"`
	Metadata         map[string]string   `json:"metadata"`
	Source           string              `json:"source"`
	Files            []*WorkflowFileInfo `json:"files"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

type WorkflowFileInfo struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type WorkflowParam struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Default  string `json:"default"`
}

type WorkflowFile struct {
	ID                string    `json:"id"`
	WorkflowVersionID string    `json:"workflowVersionID"`
	Path              string    `json:"path"`
	Content           string    `json:"content"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func WorkflowFileDTOtoWorkflowFileVO(workflowFileDTO *query.WorkflowFile) *WorkflowFile {
	if workflowFileDTO == nil {
		return nil
	}
	return &WorkflowFile{
		ID:                workflowFileDTO.ID,
		WorkflowVersionID: workflowFileDTO.WorkflowVersionID,
		Path:              workflowFileDTO.Path,
		Content:           workflowFileDTO.Content,
		CreatedAt:         workflowFileDTO.CreatedAt,
		UpdatedAt:         workflowFileDTO.UpdatedAt,
	}
}

func WorkflowVersionDTOtoVO(workflowVersion *query.WorkflowVersion) *WorkflowVersion {
	if workflowVersion == nil {
		return nil
	}
	inputs := make([]*WorkflowParam, len(workflowVersion.Inputs))
	for index, input := range workflowVersion.Inputs {
		inputs[index] = &WorkflowParam{
			Name:     input.Name,
			Type:     input.Type,
			Optional: input.Optional,
			Default:  input.Default,
		}
	}

	outputs := make([]*WorkflowParam, len(workflowVersion.Outputs))
	for index, output := range workflowVersion.Outputs {
		outputs[index] = &WorkflowParam{
			Name:     output.Name,
			Type:     output.Type,
			Optional: output.Optional,
			Default:  output.Default,
		}
	}
	// all outputs are optional
	for i := range outputs {
		outputs[i].Optional = true
	}

	files := make([]*WorkflowFileInfo, len(workflowVersion.Files))
	for index, file := range workflowVersion.Files {
		files[index] = &WorkflowFileInfo{
			ID:   file.ID,
			Path: file.Path,
		}
	}

	return &WorkflowVersion{
		ID:               workflowVersion.ID,
		Status:           workflowVersion.Status,
		Message:          workflowVersion.Message,
		Language:         workflowVersion.Language,
		LanguageVersion:  workflowVersion.LanguageVersion,
		MainWorkflowPath: workflowVersion.MainWorkflowPath,
		Graph:            workflowVersion.Graph,
		Source:           workflowVersion.Source,
		Files:            files,
		Metadata:         workflowVersion.Metadata,
		Inputs:           inputs,
		Outputs:          outputs,
		CreatedAt:        workflowVersion.CreatedAt,
		UpdatedAt:        workflowVersion.UpdatedAt,
	}
}
