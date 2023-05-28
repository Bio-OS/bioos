package convert

import (
	"reflect"

	"k8s.io/utils/pointer"

	submissionproto "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc/proto"
)

type CreateSubmissionRequest struct {
	WorkspaceID    string         `path:"workspace_id"`
	Name           string         `json:"name"`
	WorkflowID     string         `json:"workflowID"`
	Description    *string        `json:"description"`
	Type           string         `json:"type"`
	Entity         *Entity        `json:"entity"`
	ExposedOptions ExposedOptions `json:"exposedOptions"`
	InOutMaterial  *InOutMaterial `json:"inOutMaterial"`
}

func (req *CreateSubmissionRequest) ToGRPC() *submissionproto.CreateSubmissionRequest {
	return &submissionproto.CreateSubmissionRequest{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		WorkflowID:  req.WorkflowID,
		Description: pointer.StringDeref(req.Description, ""),
		Type:        req.Type,
		Entity: &submissionproto.Entity{
			DataModelID:     req.Entity.DataModelID,
			DataModelRowIDs: req.Entity.DataModelRowIDs,
			InputsTemplate:  req.Entity.InputsTemplate,
			OutputsTemplate: req.Entity.OutputsTemplate,
		},
		ExposedOptions: &submissionproto.ExposedOptions{
			ReadFromCache: req.ExposedOptions.ReadFromCache,
		},
		InOutMaterial: &submissionproto.InOutMaterial{
			InputsMaterial:  req.InOutMaterial.InputsMaterial,
			OutputsMaterial: req.InOutMaterial.OutputsMaterial,
		},
	}
}

type Entity struct {
	DataModelID     string   `json:"dataModelID"`
	DataModelRowIDs []string `json:"dataModelRowIDs"`
	/** 输入配置，json 序列化后的 string
	  采用 json 序列化原因基于以下两点考虑：
	  - thrift/接口设计层面不允许 `Value` 类型不确定
	  - 在 inputs/outputs 层级进行序列化可使得 `bioos-server` 不处理 `Inputs`/`Outputs`(非 `this.xxx` 索引的输入) 就入库/提交给计算引擎，达到透传效果
	*/
	InputsTemplate string `json:"inputsTemplate"`
	/** 输出配置，json 序列化后的 string
	  采用 json 序列化原因基于以下两点考虑：
	  - thrift/接口设计层面不允许 `Value` 类型不确定
	  - 在 inputs/outputs 层级进行序列化可使得 `bioos-server` 不处理 `Inputs`/`Outputs`(非 `this.xxx` 索引的输入) 就入库/提交给计算引擎，达到透传效果
	*/
	OutputsTemplate string `json:"outputsTemplate"`
}

type InOutMaterial struct {
	InputsMaterial  string `json:"inputsMaterial"`
	OutputsMaterial string `json:"outputsMaterial"`
}

type ExposedOptions struct {
	ReadFromCache bool `json:"readFromCache"`
}

type CreateSubmissionResponse struct {
	ID string `json:"id"`
}

func (resp *CreateSubmissionResponse) FromGRPC(protoResp *submissionproto.CreateSubmissionResponse) {
	resp.ID = protoResp.GetId()
}

type CheckSubmissionRequest struct {
	WorkspaceID string `path:"workspace_id"`
	Name        string `json:"name"`
}

func (req *CheckSubmissionRequest) ToGRPC() *submissionproto.CheckSubmissionRequest {
	return &submissionproto.CheckSubmissionRequest{}
}

type CheckSubmissionResponse struct {
	IsNameExist bool `json:"isNameExist"`
}

func (resp *CheckSubmissionResponse) FromGRPC(protoResp *submissionproto.CheckSubmissionResponse) {
	resp.IsNameExist = protoResp.GetIsNameExist()
}

type ListSubmissionsRequest struct {
	WorkspaceID string   `path:"workspace_id"`
	Page        int      `query:"page"`
	Size        int      `query:"size"`
	OrderBy     string   `query:"orderBy"`
	SearchWord  string   `query:"searchWord"`
	WorkflowID  string   `query:"workflowID,omitempty"`
	Status      []string `query:"status,omitempty"`
	IDs         []string `query:"ids,omitempty"`
}

func (req *ListSubmissionsRequest) ToGRPC() *submissionproto.ListSubmissionsRequest {
	return &submissionproto.ListSubmissionsRequest{
		WorkspaceID: req.WorkspaceID,
		Page:        int32(req.Page),
		Size:        int32(req.Size),
		OrderBy:     req.OrderBy,
		SearchWord:  req.SearchWord,
		WorkflowID:  req.WorkflowID,
		Status:      req.Status,
		Ids:         req.IDs,
	}
}

type ListSubmissionsResponse struct {
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
	Items []SubmissionItem `json:"items"`
}

type listSubmissionsResponseBriefItems struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	WorkflowID        string `json:"workflowID"`
	WorkflowVersionID string `json:"versionID"`
}

func (resp *ListSubmissionsResponse) BriefItems() reflect.Value {
	briefItems := make([]listSubmissionsResponseBriefItems, len(resp.Items))
	for i, item := range resp.Items {
		briefItems[i] = listSubmissionsResponseBriefItems{
			ID:                item.ID,
			Name:              item.Name,
			Status:            item.Status,
			WorkflowID:        item.WorkflowVersion.ID,
			WorkflowVersionID: item.WorkflowVersion.VersionID,
		}
	}
	return reflect.ValueOf(briefItems)
}

func (resp *ListSubmissionsResponse) FromGRPC(protoResp *submissionproto.ListSubmissionsResponse) {
	resp.Page = int(protoResp.GetPage())
	resp.Size = int(protoResp.GetSize())
	resp.Total = int(protoResp.GetTotal())
	resp.Items = make([]SubmissionItem, len(protoResp.GetItems()))
	for i, item := range protoResp.GetItems() {
		resp.Items[i] = SubmissionItem{
			ID:          item.GetId(),
			Name:        item.GetName(),
			Description: &item.Description,
			Status:      item.GetStatus(),
			StartTime:   item.GetStartTime(),
			FinishTime:  &item.FinishTime,
			Duration:    item.GetDuration(),
			WorkflowVersion: WorkflowVersionBrief{
				item.GetWorkflowVersion().GetId(),
				item.GetWorkflowVersion().GetVersionID(),
			},
			RunStatus: Status{
				Count:        item.GetRunStatus().GetCount(),
				Pending:      item.GetRunStatus().GetPending(),
				Succeeded:    item.GetRunStatus().GetSucceeded(),
				Failed:       item.GetRunStatus().GetFailed(),
				Running:      item.GetRunStatus().GetRunning(),
				Cancelling:   item.GetRunStatus().GetCancelling(),
				Cancelled:    item.GetRunStatus().GetCancelled(),
				Queued:       item.GetRunStatus().GetQueued(),
				Initializing: item.GetRunStatus().GetInitializing(),
			},
			Entity: &Entity{
				DataModelID:     item.GetEntity().GetDataModelID(),
				DataModelRowIDs: item.GetEntity().GetDataModelRowIDs(),
				InputsTemplate:  item.GetEntity().GetInputsTemplate(),
				OutputsTemplate: item.GetEntity().GetOutputsTemplate(),
			},
			ExposedOptions: ExposedOptions{
				ReadFromCache: item.GetExposedOptions().GetReadFromCache(),
			},
			InOutMaterial: &InOutMaterial{
				InputsMaterial:  item.GetInOutMaterial().GetInputsMaterial(),
				OutputsMaterial: item.GetInOutMaterial().GetOutputsMaterial(),
			},
		}
	}
}

type SubmissionItem struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Description     *string              `json:"description"`
	Status          string               `json:"status"`
	StartTime       int64                `json:"startTime"`
	FinishTime      *int64               `json:"finishTime"`
	Duration        int64                `json:"duration"`
	WorkflowVersion WorkflowVersionBrief `json:"workflowVersion"`
	RunStatus       Status               `json:"runStatus"`
	Entity          *Entity              `json:"entity"`
	ExposedOptions  ExposedOptions       `json:"exposedOptions"`
	InOutMaterial   *InOutMaterial       `json:"inOutMaterial"`
}

type WorkflowVersionBrief struct {
	ID        string `json:"id"`
	VersionID string `json:"versionID"`
}

type Status struct {
	Count        int64 `json:"count"`
	Pending      int64 `json:"pending"`
	Succeeded    int64 `json:"succeeded"`
	Failed       int64 `json:"failed"`
	Running      int64 `json:"running"`
	Cancelling   int64 `json:"cancelling"`
	Cancelled    int64 `json:"cancelled"`
	Queued       int64 `json:"queued"`
	Initializing int64 `json:"initializing"`
}

type DeleteSubmissionRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

func (req *DeleteSubmissionRequest) ToGRPC() *submissionproto.DeleteSubmissionRequest {
	return &submissionproto.DeleteSubmissionRequest{
		WorkspaceID: req.WorkspaceID,
		Id:          req.ID,
	}
}

type DeleteSubmissionResponse struct {
}

func (resp *DeleteSubmissionResponse) FromGRPC(protoResp *submissionproto.DeleteSubmissionResponse) {
	return
}

type CancelSubmissionRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

func (req *CancelSubmissionRequest) ToGRPC() *submissionproto.CancelSubmissionRequest {
	return &submissionproto.CancelSubmissionRequest{
		WorkspaceID: req.WorkspaceID,
		Id:          req.ID,
	}
}

type CancelSubmissionResponse struct {
}

func (resp *CancelSubmissionResponse) FromGRPC(protoResp *submissionproto.CancelSubmissionResponse) {
	return
}
