package convert

import (
	"reflect"

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

type GetDataModelRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

func (req *GetDataModelRequest) ToGRPC() *workspaceproto.GetDataModelRequest {
	return &workspaceproto.GetDataModelRequest{
		Id:          req.ID,
		WorkspaceID: req.WorkspaceID,
	}
}

type GetDataModelResponse struct {
	DataModel *DataModel `json:"dataModel,omitempty"`
	Headers   []string   `json:"headers"`
}

func (resp *GetDataModelResponse) FromGRPC(protoResp *workspaceproto.GetDataModelResponse) {
	resp.DataModel = &DataModel{
		ID:       protoResp.GetDataModel().GetId(),
		Name:     protoResp.GetDataModel().GetName(),
		RowCount: protoResp.GetDataModel().GetRowCount(),
		Type:     protoResp.GetDataModel().GetType(),
	}
	resp.Headers = protoResp.GetHeaders()
}

type ListDataModelsRequest struct {
	WorkspaceID string   `path:"workspace_id"`
	Types       []string `query:"types,omitempty"`
	SearchWord  string   `query:"searchWord"`
	IDs         []string `query:"ids,omitempty"`
}

func (req *ListDataModelsRequest) ToGRPC() *workspaceproto.ListDataModelsRequest {
	return &workspaceproto.ListDataModelsRequest{
		WorkspaceID: req.WorkspaceID,
		Types:       req.Types,
		SearchWord:  req.SearchWord,
		Ids:         req.IDs,
	}
}

type ListDataModelsResponse struct {
	Items []DataModel `json:"Items"`
}

func (resp *ListDataModelsResponse) BriefItems() interface{} {
	briefItems := make([]DataModel, len(resp.Items))
	for i, item := range resp.Items {
		briefItems[i] = item
	}
	return reflect.ValueOf(briefItems)
}

func (resp *ListDataModelsResponse) FromGRPC(protoResp *workspaceproto.ListDataModelsResponse) {
	resp.Items = make([]DataModel, len(protoResp.GetItems()))
	for i, item := range protoResp.Items {
		resp.Items[i] = DataModel{
			ID:       item.GetId(),
			Name:     item.GetName(),
			RowCount: item.GetRowCount(),
			Type:     item.GetType(),
		}
	}
}

type DataModel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RowCount int64  `json:"rowCount"`
	Type     string `json:"type"`
}

type ListDataModelRowsRequest struct {
	WorkspaceID string   `path:"workspace_id"`
	ID          string   `path:"id"`
	Page        int32    `query:"page"`
	Size        int32    `query:"size"`
	OrderBy     string   `query:"orderBy"`
	SearchWord  string   `query:"searchWord"`
	InSetIDs    []string `query:"inSetIDs"`
	RowIDs      []string `query:"rowIDs"`
}

func (req *ListDataModelRowsRequest) ToGRPC() *workspaceproto.ListDataModelRowsRequest {
	return &workspaceproto.ListDataModelRowsRequest{
		WorkspaceID: req.WorkspaceID,
		Id:          req.ID,
		Page:        req.Page,
		Size:        req.Size,
		OrderBy:     req.OrderBy,
		SearchWord:  req.SearchWord,
		InSetIDs:    req.InSetIDs,
		RowIDs:      req.RowIDs,
	}
}

type ListDataModelRowsResponse struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
	Page    int32      `json:"page"`
	Size    int32      `json:"size"`
	Total   int64      `json:"total"`
}

func (resp *ListDataModelRowsResponse) FromGRPC(protoResp *workspaceproto.ListDataModelRowsResponse) {
	resp.Headers = protoResp.GetHeaders()
	resp.Rows = make([][]string, len(protoResp.Rows))
	for i, r := range protoResp.Rows {
		resp.Rows[i] = r.Grids
	}
	resp.Page = protoResp.GetPage()
	resp.Size = protoResp.GetSize()
	resp.Total = protoResp.GetTotal()
}

type Row struct {
	Grid []Grid `json:"grid"`
}

type Grid struct {
	Value []byte `json:"value"`
	Type  string `json:"type"`
}

type PatchDataModelRequest struct {
	WorkspaceID string     `path:"workspace_id"`
	Name        string     `json:"name"`
	Async       bool       `json:"async"`
	Headers     []string   `json:"headers"`
	Rows        [][]string `json:"rows"`
}

func (req *PatchDataModelRequest) ToGRPC() *workspaceproto.PatchDataModelRequest {
	out := &workspaceproto.PatchDataModelRequest{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Async:       req.Async,
		Headers:     req.Headers,
	}
	out.Rows = make([]*workspaceproto.Row, len(req.Rows))
	for i, r := range req.Rows {
		out.Rows[i] = &workspaceproto.Row{
			Grids: r,
		}
	}
	return out
}

type PatchDataModelResponse struct {
	ID string `json:"id"`
}

func (resp *PatchDataModelResponse) FromGRPC(protoResp *workspaceproto.PatchDataModelResponse) {
	resp.ID = protoResp.GetId()
}

type DeleteDataModelRequest struct {
	ID          string   `path:"id"`
	WorkspaceID string   `path:"workspace_id"`
	Headers     []string `query:"headers,omitempty"`
	RowIDs      []string `query:"rowIDs,omitempty"`
}

func (req *DeleteDataModelRequest) ToGRPC() *workspaceproto.DeleteDataModelRequest {
	return &workspaceproto.DeleteDataModelRequest{
		Id:          req.ID,
		WorkspaceID: req.WorkspaceID,
		Headers:     req.Headers,
		RowIDs:      req.RowIDs,
	}
}

type DeleteDataModelResponse struct{}

func (resp *DeleteDataModelResponse) FromGRPC(protoResp *workspaceproto.DeleteDataModelResponse) {
	return
}

type ListAllDataModelRowIDsRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

func (req *ListAllDataModelRowIDsRequest) ToGRPC() *workspaceproto.ListAllDataModelRowIDsRequest {
	return &workspaceproto.ListAllDataModelRowIDsRequest{
		WorkspaceID: req.WorkspaceID,
		Id:          req.ID,
	}
}

type ListAllDataModelRowIDsResponse struct {
	RowIDs []string `json:"rowIDs"`
}

func (resp *ListAllDataModelRowIDsResponse) FromGRPC(protoResp *workspaceproto.ListAllDataModelRowIDsResponse) {
	resp.RowIDs = protoResp.GetRowIDs()
}
