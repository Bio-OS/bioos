package mysql

import (
	"context"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
)

func WorkspacePOToWorkspaceDTO(ctx context.Context, w *Workspace) *query.WorkspaceItem {
	item := &query.WorkspaceItem{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		CreatedAt:   w.CreateTime,
		UpdatedAt:   w.UpdateTime,
	}
	if w.Storage.NFS != nil {
		item.Storage = query.WorkspaceStorage{NFS: &query.NFSWorkspaceStorage{MountPath: w.Storage.NFS.MountPath}}
	}
	return item
}

func WorkspacePOToWorkspaceDO(ctx context.Context, w *Workspace) (*workspace.Workspace, error) {
	factory := workspace.NewWorkspaceFactory(ctx)
	param := workspace.CreateWorkspaceParam{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		CreatedAt:   w.CreateTime,
		UpdatedAt:   w.UpdateTime,
	}
	if w.Storage.NFS != nil {
		param.Storage.NFS = &workspace.NFSStorage{MountPath: w.Storage.NFS.MountPath}
	}
	return factory.CreateWithWorkspaceParam(param)
}

func WorkspaceDOtoWorkspacePO(w *workspace.Workspace) *Workspace {
	res := &Workspace{
		ID:          w.GetID(),
		Name:        w.GetName(),
		Description: w.GetDescription(),
		CreateTime:  w.GetCreatedAt(),
		UpdateTime:  w.GetUpdatedAt(),
	}
	storage := w.GetStorage()
	if storage.NFS != nil {
		res.Storage = WorkspaceStorage{
			NFS: &NFSWorkspaceStorage{MountPath: storage.NFS.MountPath},
		}
	}
	return res
}
