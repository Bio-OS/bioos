package handlers

type GetDataModelRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

type GetDataModelResponse struct {
	DataModel *DataModel `json:"dataModel,omitempty"`
	Headers   []string   `json:"headers"`
}

type ListDataModelsRequest struct {
	WorkspaceID string   `path:"workspace_id"`
	Types       []string `query:"types"`
	SearchWord  string   `query:"searchWord"`
	IDs         []string `query:"ids"`
}

type ListDataModelsResponse struct {
	Items []*DataModel `json:"Items"`
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

type ListDataModelRowsResponse struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
	Page    int32      `json:"page"`
	Size    int32      `json:"size"`
	Total   int64      `json:"total"`
}

type Row struct {
	Grid []*Grid `json:"grid"`
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

type PatchDataModelResponse struct {
	ID string `json:"id"`
}

type DeleteDataModelRequest struct {
	ID          string   `path:"id"`
	WorkspaceID string   `path:"workspace_id"`
	Headers     []string `query:"headers"`
	RowIDs      []string `query:"rowIDs"`
}

type ListAllDataModelRowIDsRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

type ListAllDataModelRowIDsResponse struct {
	RowIDs []string `json:"rowIDs"`
}
