package convert

import (
	"reflect"
	"strings"

	"k8s.io/utils/pointer"

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
)

type CreateWorkspaceRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Storage     *WorkspaceStorage `json:"storage"`
}

func (req *CreateWorkspaceRequest) ToGRPC() *workspaceproto.CreateWorkspaceRequest {
	return &workspaceproto.CreateWorkspaceRequest{
		Name:        req.Name,
		Description: req.Description,
		Storage: &workspaceproto.WorkspaceStorage{
			Nfs: &workspaceproto.NFSWorkspaceStorage{
				MountPath: req.Storage.NFS.MountPath,
			},
		},
	}
}

type CreateWorkspaceResponse struct {
	Id string `json:"id"`
}

func (resp *CreateWorkspaceResponse) FromGRPC(protoResp *workspaceproto.CreateWorkspaceResponse) {
	resp.Id = protoResp.Id
}

type DeleteWorkspaceRequest struct {
	Id string `path:"id"`
}

func (req *DeleteWorkspaceRequest) ToGRPC() *workspaceproto.DeleteWorkspaceRequest {
	return &workspaceproto.DeleteWorkspaceRequest{
		Id: req.Id,
	}
}

type DeleteWorkspaceResponse struct {
}

func (resp *DeleteWorkspaceResponse) FromGRPC(protoResp *workspaceproto.DeleteWorkspaceResponse) {
	return
}

type UpdateWorkspaceRequest struct {
	ID          string  `path:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (req *UpdateWorkspaceRequest) ToGRPC() *workspaceproto.UpdateWorkspaceRequest {
	return &workspaceproto.UpdateWorkspaceRequest{
		Id:          req.ID,
		Name:        pointer.StringDeref(req.Name, ""),
		Description: pointer.StringDeref(req.Description, ""),
	}
}

type UpdateWorkspaceResponse struct {
}

func (resp *UpdateWorkspaceResponse) FromGRPC(protoResp *workspaceproto.UpdateWorkspaceResponse) {
	return
}

type GetWorkspaceRequest struct {
	Id string `path:"id"`
}

func (req *GetWorkspaceRequest) ToGRPC() *workspaceproto.GetWorkspaceRequest {
	r := &workspaceproto.GetWorkspaceRequest{
		Id: req.Id,
	}
	return r
}

type GetWorkspaceResponse struct {
	WorkspaceItem `json:",inline"`
}

func (resp *GetWorkspaceResponse) FromGRPC(protoResp *workspaceproto.GetWorkspaceResponse) {
	resp.WorkspaceItem = WorkspaceItem{
		Id:          protoResp.Workspace.Id,
		Name:        protoResp.Workspace.Name,
		Description: protoResp.Workspace.Description,
		Storage: &WorkspaceStorage{
			NFS: &NFSWorkspaceStorage{
				MountPath: protoResp.Workspace.Storage.Nfs.MountPath,
			},
		},
		CreateTime: protoResp.Workspace.CreatedAt.GetSeconds(),
		UpdateTime: protoResp.Workspace.UpdatedAt.GetSeconds(),
	}
	return
}

type ListWorkspacesRequest struct {
	Page       int    `query:"page"`
	Size       int    `query:"size"`
	OrderBy    string `query:"orderBy"`
	SearchWord string `query:"searchWord"`
	IDs        string `query:"ids"`
}

func (req *ListWorkspacesRequest) ToGRPC() *workspaceproto.ListWorkspaceRequest {
	r := &workspaceproto.ListWorkspaceRequest{
		Page:       int32(req.Page),
		Size:       int32(req.Size),
		OrderBy:    req.OrderBy,
		SearchWord: req.SearchWord,
	}
	if req.IDs != "" {
		r.Ids = strings.Split(req.IDs, consts.QuerySliceDelimiter)
	}
	return r
}

type ListWorkspacesResponse struct {
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	Items []WorkspaceItem `json:"items"`
}

type listWorkspacesBriefItems struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MountType   string `json:"mount_type"`
	MountPath   string `json:"mount_path"`
}

func (resp *ListWorkspacesResponse) BriefItems() reflect.Value {
	briefItems := make([]listWorkspacesBriefItems, len(resp.Items))
	for i, item := range resp.Items {
		briefItems[i] = listWorkspacesBriefItems{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			MountType:   "nfs",
			MountPath:   item.Storage.NFS.MountPath,
		}
	}
	return reflect.ValueOf(briefItems)
}

func (resp *ListWorkspacesResponse) FromGRPC(protoResp *workspaceproto.ListWorkspaceResponse) {
	resp.Page = int(protoResp.Page)
	resp.Size = int(protoResp.Size)
	resp.Total = int(protoResp.Total)
	resp.Items = make([]WorkspaceItem, len(protoResp.GetItems()))
	for i, item := range protoResp.Items {
		resp.Items[i] = WorkspaceItem{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Storage: &WorkspaceStorage{
				NFS: &NFSWorkspaceStorage{
					MountPath: item.Storage.Nfs.MountPath,
				},
			},
			CreateTime: item.CreatedAt.GetSeconds(),
			UpdateTime: item.UpdatedAt.GetSeconds(),
		}
	}
	return
}

type WorkspaceItem struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Storage     *WorkspaceStorage `json:"storage"`
	CreateTime  int64             `json:"createTime"`
	UpdateTime  int64             `json:"updateTime"`
}

type WorkspaceStorage struct {
	NFS *NFSWorkspaceStorage `json:"nfs,omitempty"`
}

// For table/text output
func (s WorkspaceStorage) String() string {
	return s.NFS.MountPath
}

type NFSWorkspaceStorage struct {
	MountPath string `json:"mountPath"`
}

type ImportWorkspaceRequest struct {
	FilePath  string
	MountPath string `query:"mountPath"`
	MountType string `query:"mountType"`
}

type ImportWorkspaceResponse struct {
	Id string `json:"id"`
}

func (resp *ImportWorkspaceResponse) FromGRPC(protoResp *workspaceproto.ImportWorkspaceResponse) {
	resp.Id = protoResp.GetId()
}
