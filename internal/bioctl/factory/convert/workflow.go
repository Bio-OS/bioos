package convert

import (
	"reflect"
	"strings"
	"time"

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
)

type CreateWorkflowRequest struct {
	ID               string `json:"id"`
	WorkspaceID      string `path:"workspace-id"`
	Name             string `json:"name" validate:"required,resName"`
	Description      string `json:"description" validate:"required,workspaceDesc"`
	Language         string `json:"language" validate:"required,oneof=WDL"`
	Source           string `json:"source" validate:"required,oneof=git"`
	Url              string `json:"url" validate:"required"`
	Tag              string `json:"tag" validate:"required"`
	Token            string `json:"token"`
	MainWorkflowPath string `json:"mainWorkflowPath" validate:"required"`
}

func (req *CreateWorkflowRequest) ToGRPC() *workspaceproto.CreateWorkflowRequest {
	return &workspaceproto.CreateWorkflowRequest{
		Id:               &req.ID,
		WorkspaceID:      req.WorkspaceID,
		Name:             req.Name,
		Description:      &req.Description,
		Language:         req.Language,
		Source:           req.Source,
		Url:              &req.Url,
		Tag:              &req.Tag,
		Token:            &req.Token,
		MainWorkflowPath: req.MainWorkflowPath,
	}
}

type CreateWorkflowResponse struct {
	ID string `json:"id"`
}

func (resp *CreateWorkflowResponse) FromGRPC(protoResp *workspaceproto.CreateWorkflowResponse) {
	resp.ID = protoResp.GetId()
}

type GetWorkflowRequest struct {
	WorkspaceID string `path:"workspace-id"`
	ID          string `path:"id"`
}

func (req *GetWorkflowRequest) ToGRPC() *workspaceproto.GetWorkflowRequest {
	return &workspaceproto.GetWorkflowRequest{
		WorkspaceID: req.WorkspaceID,
		Id:          req.ID,
	}
}

type GetWorkflowResponse struct {
	Workflow WorkflowItem `json:"workflow"`
}

func (resp *GetWorkflowResponse) FromGRPC(protoResp *workspaceproto.GetWorkflowResponse) {
	resp.Workflow = convertWorkflowItem(protoResp.GetWorkflow())
}

type ListWorkflowsRequest struct {
	Page        int    `query:"page"`
	Size        int    `query:"size"`
	OrderBy     string `query:"orderBy"`
	SearchWord  string `query:"searchWord"`
	IDs         string `query:"ids"`
	WorkspaceID string `path:"workspace-id"`
}

func (req *ListWorkflowsRequest) ToGRPC() *workspaceproto.ListWorkflowRequest {
	return &workspaceproto.ListWorkflowRequest{
		Page:        int32(req.Page),
		Size:        int32(req.Size),
		OrderBy:     req.OrderBy,
		SearchWord:  req.SearchWord,
		Ids:         strings.Split(req.IDs, consts.QuerySliceDelimiter),
		WorkspaceID: req.WorkspaceID,
	}
}

type ListWorkflowsResponse struct {
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Total int            `json:"total"`
	Items []WorkflowItem `json:"items"`
}

type listWorkflowsResponseBriefItems struct {
	ID                            string `json:"id"`
	Name                          string `json:"name"`
	Description                   string `json:"description"`
	LatestVersionStatus           string `json:"status"`
	LatestVersionMessage          string `json:"message"`
	LatestVersionLanguage         string `json:"language"`
	LatestVersionLanguageVersion  string `json:"languageVersion"`
	LatestVersionMainWorkflowPath string `json:"mainWorkflowPath"`
	LatestVersionSource           string `json:"source"`
}

func (resp *ListWorkflowsResponse) BriefItems() reflect.Value {
	briefItems := make([]listWorkflowsResponseBriefItems, len(resp.Items))
	for i, item := range resp.Items {
		briefItems[i] = listWorkflowsResponseBriefItems{
			ID:                            item.ID,
			Name:                          item.Name,
			Description:                   item.Description,
			LatestVersionStatus:           item.LatestVersion.Status,
			LatestVersionMessage:          item.LatestVersion.Message,
			LatestVersionLanguage:         item.LatestVersion.Language,
			LatestVersionLanguageVersion:  item.LatestVersion.LanguageVersion,
			LatestVersionMainWorkflowPath: item.LatestVersion.MainWorkflowPath,
			LatestVersionSource:           item.LatestVersion.Source,
		}
	}
	return reflect.ValueOf(briefItems)
}

func (resp *ListWorkflowsResponse) FromGRPC(protoResp *workspaceproto.ListWorkflowResponse) {
	resp.Page = int(protoResp.Page)
	resp.Size = int(protoResp.Size)
	resp.Total = int(protoResp.Total)
	resp.Items = make([]WorkflowItem, len(protoResp.Items))
	for i := range protoResp.Items {
		resp.Items[i] = convertWorkflowItem(protoResp.Items[i])
	}

}

type UpdateWorkflowRequest struct {
	ID          string `path:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WorkspaceID string `path:"workspace-id"`
}

func (req *UpdateWorkflowRequest) ToGRPC() *workspaceproto.UpdateWorkflowRequest {
	return &workspaceproto.UpdateWorkflowRequest{
		Id:          req.ID,
		Name:        &req.Name,
		Description: &req.Description,
		WorkspaceID: req.WorkspaceID,
	}
}

type UpdateWorkflowResponse struct {
}

func (resp *UpdateWorkflowResponse) FromGRPC(protoResp *workspaceproto.UpdateWorkflowResponse) {
	return
}

type DeleteWorkflowRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
}

func (req *DeleteWorkflowRequest) ToGRPC() *workspaceproto.DeleteWorkflowRequest {
	return &workspaceproto.DeleteWorkflowRequest{
		Id:          req.ID,
		WorkspaceID: req.WorkspaceID,
	}
}

type DeleteWorkflowResponse struct {
}

func (resp *DeleteWorkflowResponse) FromGRPC(protoResp *workspaceproto.DeleteWorkflowResponse) {
	return
}

type GetWorkflowFileRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
	WorkflowID  string `path:"workflow-id"`
}

func (req *GetWorkflowFileRequest) ToGRPC() *workspaceproto.GetWorkflowFileRequest {
	return &workspaceproto.GetWorkflowFileRequest{
		Id:          req.ID,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
	}
}

type GetWorkflowFileResponse struct {
	File *WorkflowFile `json:"file"`
}

func (resp *GetWorkflowFileResponse) FromGRPC(protoResp *workspaceproto.GetWorkflowFileResponse) {
	resp.File = &WorkflowFile{
		ID:                protoResp.GetFile().GetId(),
		WorkflowVersionID: protoResp.GetFile().GetWorkflowVersionID(),
		Path:              protoResp.GetFile().GetPath(),
		Content:           protoResp.GetFile().GetContent(),
		CreatedAt:         protoResp.GetFile().GetCreatedAt().AsTime(),
		UpdatedAt:         protoResp.GetFile().GetUpdatedAt().AsTime(),
	}
}

type GetWorkflowVersionRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
	WorkflowID  string `path:"workflow-id"`
}

func (req *GetWorkflowVersionRequest) ToGRPC() *workspaceproto.GetWorkflowVersionRequest {
	return &workspaceproto.GetWorkflowVersionRequest{
		Id:          req.ID,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
	}
}

type GetWorkflowVersionResponse struct {
	Version WorkflowVersion `json:"version"`
}

func (resp *GetWorkflowVersionResponse) FromGRPC(protoResp *workspaceproto.GetWorkflowVersionResponse) {
	resp.Version = convertWorkflowVersion(protoResp.GetVersion())
}

type ListWorkflowFilesRequest struct {
	Page              int    `query:"page"`
	Size              int    `query:"size"`
	OrderBy           string `query:"orderBy"`
	IDs               string `query:"ids"`
	WorkspaceID       string `path:"workspace-id"`
	WorkflowID        string `path:"workflow-id"`
	WorkflowVersionID string `query:"workflowVersionID,omitempty"`
}

func (req *ListWorkflowFilesRequest) ToGRPC() *workspaceproto.ListWorkflowFilesRequest {
	return &workspaceproto.ListWorkflowFilesRequest{
		Page:              int32(req.Page),
		Size:              int32(req.Size),
		OrderBy:           req.OrderBy,
		Ids:               strings.Split(req.IDs, consts.QuerySliceDelimiter),
		WorkspaceID:       req.WorkspaceID,
		WorkflowID:        req.WorkflowID,
		WorkflowVersionID: &req.WorkflowVersionID,
	}
}

type ListWorkflowFilesResponse struct {
	Page        int             `json:"page"`
	Size        int             `json:"size"`
	Total       int             `json:"total"`
	WorkspaceID string          `json:"workspaceID"`
	WorkflowID  string          `json:"workflowID"`
	Items       []*WorkflowFile `json:"items"`
}

func (resp *ListWorkflowFilesResponse) FromGRPC(protoResp *workspaceproto.ListWorkflowFilesResponse) {
	resp.Page = int(protoResp.GetPage())
	resp.Size = int(protoResp.GetSize())
	resp.Total = int(protoResp.GetTotal())
	resp.WorkspaceID = protoResp.GetWorkspaceID()
	resp.WorkflowID = protoResp.GetWorkflowID()
	resp.Items = make([]*WorkflowFile, len(protoResp.GetFiles()))
	for i, f := range protoResp.GetFiles() {
		resp.Items[i] = &WorkflowFile{
			ID:                f.GetId(),
			WorkflowVersionID: f.GetWorkflowVersionID(),
			Path:              f.GetPath(),
			Content:           f.GetContent(),
			CreatedAt:         f.GetCreatedAt().AsTime(),
			UpdatedAt:         f.GetUpdatedAt().AsTime(),
		}
	}

}

type ListWorkflowVersionsRequest struct {
	Page        int    `query:"page"`
	Size        int    `query:"size"`
	OrderBy     string `query:"orderBy"`
	IDs         string `query:"ids"`
	WorkspaceID string `path:"workspace-id"`
	WorkflowID  string `path:"workflow-id"`
}

func (req *ListWorkflowVersionsRequest) ToGRPC() *workspaceproto.ListWorkflowVersionsRequest {
	return &workspaceproto.ListWorkflowVersionsRequest{
		Page:        int32(req.Page),
		Size:        int32(req.Size),
		OrderBy:     req.OrderBy,
		Ids:         strings.Split(req.IDs, consts.QuerySliceDelimiter),
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
	}
}

type ListWorkflowVersionsResponse struct {
	Page        int               `json:"page"`
	Size        int               `json:"size"`
	Total       int               `json:"total"`
	WorkspaceID string            `json:"workspaceID"`
	WorkflowID  string            `json:"workflowID"`
	Items       []WorkflowVersion `json:"items"`
}

func (resp *ListWorkflowVersionsResponse) FromGRPC(protoResp *workspaceproto.ListWorkflowVersionsResponse) {
	resp.Page = int(protoResp.GetPage())
	resp.Size = int(protoResp.GetSize())
	resp.Total = int(protoResp.GetTotal())
	resp.WorkspaceID = protoResp.GetWorkspaceID()
	resp.WorkflowID = protoResp.GetWorkflowID()
	resp.Items = make([]WorkflowVersion, len(protoResp.GetItems()))
	for i, item := range protoResp.GetItems() {
		resp.Items[i] = convertWorkflowVersion(item)
	}

}

type WorkflowItem struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	LatestVersion WorkflowVersion `json:"latestVersion"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"UpdatedAt"`
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
	UpdatedAt        time.Time           `json:"UpdatedAt"`
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
	UpdatedAt         time.Time `json:"UpdatedAt"`
}

func convertWorkflowItem(workflow *workspaceproto.Workflow) WorkflowItem {
	item := WorkflowItem{
		ID:            workflow.GetId(),
		Name:          workflow.GetName(),
		Description:   workflow.GetDescription(),
		LatestVersion: convertWorkflowVersion(workflow.GetLatestVersion()),
		CreatedAt:     workflow.GetCreatedAt().AsTime(),
		UpdatedAt:     workflow.GetUpdatedAt().AsTime(),
	}
	return item
}

func convertWorkflowVersion(version *workspaceproto.WorkflowVersion) WorkflowVersion {
	v := WorkflowVersion{
		ID:               version.GetId(),
		Status:           version.GetStatus(),
		Message:          version.GetMessage(),
		Language:         version.GetLanguage(),
		LanguageVersion:  version.GetLanguageVersion(),
		MainWorkflowPath: version.GetMainWorkflowPath(),
		Inputs:           make([]*WorkflowParam, len(version.GetInputs())),
		Outputs:          make([]*WorkflowParam, len(version.GetOutputs())),
		Graph:            version.GetGraph(),
		Metadata:         version.GetMetadata(),
		Source:           version.GetSource(),
		Files:            make([]*WorkflowFileInfo, len(version.GetFiles())),
		CreatedAt:        version.GetCreatedAt().AsTime(),
		UpdatedAt:        version.GetUpdatedAt().AsTime(),
	}
	for i, param := range version.GetInputs() {
		v.Inputs[i] = &WorkflowParam{
			Name:     param.Name,
			Type:     param.Type,
			Optional: param.Optional,
			Default:  *param.Default,
		}
	}
	for i, param := range version.GetOutputs() {
		v.Outputs[i] = &WorkflowParam{
			Name:     param.Name,
			Type:     param.Type,
			Optional: param.Optional,
			Default:  *param.Default,
		}
	}
	for i, f := range version.GetFiles() {
		v.Files[i] = &WorkflowFileInfo{
			ID:   f.Id,
			Path: f.Path,
		}
	}
	return v
}
