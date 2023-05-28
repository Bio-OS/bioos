package hertz

import (
	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/command"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

type createRequest struct {
	WorkspaceID  string                `json:"-" path:"workspace-id"`
	Image        string                `json:"image"`
	ResourceSize notebook.ResourceSize `json:"resourceSize"`
}

func (req *createRequest) toDTO() *command.CreateCommand {
	return &command.CreateCommand{
		WorkspaceID:  req.WorkspaceID,
		Image:        req.Image,
		ResourceSize: req.ResourceSize,
	}
}

type createResponse struct {
	ID string `json:"id"`
}

type updateSettingsRequest struct {
	ID           string                 `json:"-" path:"id"`
	WorkspaceID  string                 `json:"-" path:"workspace-id"`
	Image        *string                `json:"image"`
	ResourceSize *notebook.ResourceSize `json:"resourceSize"`
}

func (r *updateSettingsRequest) toDTO() *command.UpdateCommand {
	return &command.UpdateCommand{
		ID:           r.ID,
		WorkspaceID:  r.WorkspaceID,
		Image:        r.Image,
		ResourceSize: r.ResourceSize,
	}
}

type switchRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
	On          string `query:"on"` // bool will cause binding error: parameter type does not match binding data
	Off         string `query:"off"`
}

func (req *switchRequest) toDTO(onoff bool) *command.SwitchCommand {
	return &command.SwitchCommand{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
		OnOff:       onoff,
	}
}

type deleteRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
}

func (req *deleteRequest) toDTO() *command.DeleteCommand {
	return &command.DeleteCommand{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
	}
}

type listRequest struct {
	WorkspaceID string `path:"workspace-id"`
}

func (req *listRequest) toDTO() *query.ListQuery {
	return &query.ListQuery{
		WorkspaceID: req.WorkspaceID,
	}
}

type listResponseItem struct {
	ID           string                `json:"id"`
	Image        string                `json:"image"`
	ResourceSize notebook.ResourceSize `json:"resourceSize"`
	Status       string                `json:"status"`
	CreateTime   int64                 `json:"createTime"`
	UpdateTime   int64                 `json:"updateTime"`
}

func newListResponseItem(srv *query.NotebookServer) *listResponseItem {
	return &listResponseItem{
		ID:           srv.ID,
		Image:        srv.Image,
		ResourceSize: srv.ResourceSize,
		Status:       srv.Status,
		CreateTime:   srv.CreateTime.Unix(),
		UpdateTime:   srv.UpdateTime.Unix(),
	}
}

type getRequest struct {
	ID          string `path:"id"`
	WorkspaceID string `path:"workspace-id"`
	Notebook    string `query:"notebook"`
}

func (req *getRequest) toDTO() *query.GetQuery {
	return &query.GetQuery{
		ID:           req.ID,
		WorkspaceID:  req.WorkspaceID,
		EditNotebook: req.Notebook,
	}
}

type getResponse struct {
	ID           string                `json:"id"`
	Image        string                `json:"image"`
	ResourceSize notebook.ResourceSize `json:"resourceSize"`
	Status       string                `json:"status"`
	AccessURL    string                `json:"accessURL"`
	CreateTime   int64                 `json:"createTime"`
	UpdateTime   int64                 `json:"updateTime"`
}

func newGetResponse(srv *query.NotebookServer) *getResponse {
	return &getResponse{
		ID:           srv.ID,
		Image:        srv.Image,
		ResourceSize: srv.ResourceSize,
		Status:       srv.Status,
		AccessURL:    srv.AccessURL,
		CreateTime:   srv.CreateTime.Unix(),
		UpdateTime:   srv.UpdateTime.Unix(),
	}
}
