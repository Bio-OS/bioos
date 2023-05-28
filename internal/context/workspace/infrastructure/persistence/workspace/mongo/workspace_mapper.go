package mongo

import (
	"context"
	"time"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
)

func workspacePOToWorkspaceDO(ctx context.Context, w *workspacePO) (*workspace.Workspace, error) {
	factory := workspace.NewWorkspaceFactory(ctx)
	param := workspace.CreateWorkspaceParam{
		Name:        w.Name,
		Description: w.Description,
		CreatedAt:   w.CreateTime,
	}
	if w.Storage.NFS != nil {
		param.Storage.NFS = &workspace.NFSStorage{MountPath: w.Storage.NFS.MountPath}
	}
	return factory.CreateWithWorkspaceParam(param)
}

func workspaceDOtoWorkspacePO(ctx context.Context, w *workspace.Workspace) (*workspacePO, error) {
	res := &workspacePO{
		ID:          w.GetID(),
		Name:        w.GetName(),
		Description: w.GetDescription(),
		CreateTime:  w.GetCreatedAt(),
		UpdateTime:  time.Now(),
	}
	storage := w.GetStorage()
	if storage.NFS != nil {
		res.Storage = workspaceStorage{
			NFS: &workspaceStorageNFS{MountPath: storage.NFS.MountPath},
		}
	}
	return res, nil
}

func workspacePOToQueryItem(ctx context.Context, w *workspacePO) (*query.WorkspaceItem, error) {
	res := &query.WorkspaceItem{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		CreatedAt:   w.CreateTime,
		UpdatedAt:   w.UpdateTime,
	}
	if w.Storage.NFS != nil {
		res.Storage.NFS = &query.NFSWorkspaceStorage{MountPath: w.Storage.NFS.MountPath}
	}
	return res, nil
}
