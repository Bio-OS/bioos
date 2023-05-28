package handlers

type CreateWorkspaceRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Storage     WorkspaceStorage `json:"storage"`
}

type CreateWorkspaceResponse struct {
	Id string `json:"id"`
}

type DeleteWorkspaceRequest struct {
	Id string `path:"id"`
}

type UpdateWorkspaceRequest struct {
	ID          string  `path:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type GetWorkspaceByIdRequest struct {
	Id string `path:"id"`
}

type GetWorkspaceByIdResponse struct {
	WorkspaceItem `json:",inline"`
}

type ListWorkspacesRequest struct {
	Page       int    `query:"page"`
	Size       int    `query:"size"`
	OrderBy    string `query:"orderBy"`
	SearchWord string `query:"searchWord"`
	Exact      bool   `query:"exact"`
	IDs        string `query:"ids"`
}

type ListWorkspacesResponse struct {
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	Items []WorkspaceItem `json:"items"`
}

type WorkspaceItem struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Storage     WorkspaceStorage `json:"storage"`
	CreateTime  int64            `json:"createTime"`
	UpdateTime  int64            `json:"updateTime"`
}

type WorkspaceStorage struct {
	NFS *NFSWorkspaceStorage `json:"nfs,omitempty"`
}

type NFSWorkspaceStorage struct {
	MountPath string `json:"mountPath"`
}

type ImportWorkspaceRequest struct {
	FileName  string
	MountPath string `query:"mountPath"`
	MountType string `query:"mountType"`
}

type ImportWorkspaceResponse struct {
	Id string `json:"id"`
}
