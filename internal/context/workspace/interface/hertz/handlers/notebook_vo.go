package handlers

import (
	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/notebook"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
)

type createNotebookRequest struct {
	WorkspaceID string `path:"workspace-id"`
	Name        string `path:"name"`
	Content     []byte `raw_body:"required"`
}

func (req *createNotebookRequest) toDTO() *command.CreateCommand {
	return &command.CreateCommand{
		Name:        req.Name,
		WorkspaceID: req.WorkspaceID,
		Content:     req.Content,
	}
}

type getNotebookRequest struct {
	WorkspaceID string `path:"workspace-id"`
	Name        string `path:"name"`
}

func (req *getNotebookRequest) toDTO() *query.GetQuery {
	return &query.GetQuery{
		Name:        req.Name,
		WorkspaceID: req.WorkspaceID,
	}
}

type listNotebookRequest struct {
	WorkspaceID string `path:"workspace-id"`
}

func (req *listNotebookRequest) toDTO() *query.ListQuery {
	return &query.ListQuery{
		WorkspaceID: req.WorkspaceID,
	}
}

type listNotebooksResponse struct {
	Items []*notebookItem `json:"items"`
}

func newListNotebooksResponse(list []*query.Notebook) *listNotebooksResponse {
	res := &listNotebooksResponse{
		Items: make([]*notebookItem, len(list)),
	}
	for i := range list {
		res.Items[i] = newNotebookItem(list[i])
	}
	return res
}

type notebookItem struct {
	Name          string `json:"name"`
	ContentLength int64  `json:"contentLength"`
	UpdateTime    int64  `json:"updateTime"`
}

func newNotebookItem(dto *query.Notebook) *notebookItem {
	return &notebookItem{
		Name:          dto.Name,
		ContentLength: dto.Size,
		UpdateTime:    dto.UpdateTime.Unix(),
	}
}

type deleteNotebookRequest struct {
	WorkspaceID string `path:"workspace-id"`
	Name        string `path:"name"`
}

func (req *deleteNotebookRequest) toDTO() *command.DeleteCommand {
	return &command.DeleteCommand{
		Name:        req.Name,
		WorkspaceID: req.WorkspaceID,
	}
}
