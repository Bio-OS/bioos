package handlers

import (
	"strings"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workspace"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func createWorkspaceVoToDto(req CreateWorkspaceRequest) *command.CreateWorkspaceCommand {
	return &command.CreateWorkspaceCommand{
		Name:        req.Name,
		Description: req.Description,
		Storage:     workspaceStorageVoToDto(req.Storage),
	}
}

func importWorkspaceVoToDto(req ImportWorkspaceRequest) *command.ImportWorkspaceCommand {
	return &command.ImportWorkspaceCommand{
		ID:       utils.GenWorkspaceID(),
		FileName: req.FileName,
		Storage: command.WorkspaceStorage{
			NFS: &command.NFSWorkspaceStorage{
				MountPath: req.MountPath,
			},
		},
	}
}

func workspaceStorageVoToDto(s WorkspaceStorage) (res command.WorkspaceStorage) {
	if s.NFS != nil {
		res.NFS = &command.NFSWorkspaceStorage{MountPath: s.NFS.MountPath}
	}
	return res
}

func updateWorkspaceVoToDto(req UpdateWorkspaceRequest) *command.UpdateWorkspaceCommand {
	return &command.UpdateWorkspaceCommand{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}
}

func listWorkspacesVoToDto(req ListWorkspacesRequest) (*query.ListWorkspacesQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}
	filter := &query.ListWorkspacesFilter{}
	if req.SearchWord != "" {
		filter.SearchWord = req.SearchWord
		filter.Exact = req.Exact
	}
	if len(req.IDs) != 0 {
		filter.IDs = strings.Split(req.IDs, consts.QuerySliceDelimiter)
	}
	return &query.ListWorkspacesQuery{
		Pg:     *pg,
		Filter: filter,
	}, nil
}

func workspaceItemDtoToVo(ws *query.WorkspaceItem) WorkspaceItem {
	return WorkspaceItem{
		Id:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description,
		Storage:     workspaceStorageDtoToVo(ws.Storage),
		CreateTime:  ws.CreatedAt.Unix(),
		UpdateTime:  ws.UpdatedAt.Unix(),
	}
}

func workspaceStorageDtoToVo(s query.WorkspaceStorage) (res WorkspaceStorage) {
	if s.NFS != nil {
		res.NFS = &NFSWorkspaceStorage{MountPath: s.NFS.MountPath}
	}
	return res
}
