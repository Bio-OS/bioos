package command

import (
	"context"
	"fmt"
	"path"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/storage"
)

type CreateCommand struct {
	WorkspaceID  string
	Image        string
	ResourceSize notebook.ResourceSize
}

type CreateHandler interface {
	Handle(context.Context, *CreateCommand) (string, error)
}

func NewCreateHandler(svc domain.Service, factory *domain.Factory, workspaceService proto.WorkspaceServiceServer, storageOpts *storage.Options) CreateHandler {
	return &createHandler{
		factory:         factory,
		service:         svc,
		workspaceClient: workspaceService,
		storageOpts:     storageOpts,
	}
}

type createHandler struct {
	factory         *domain.Factory
	service         domain.Service
	workspaceClient proto.WorkspaceServiceServer
	storageOpts     *storage.Options
}

func (h *createHandler) Handle(ctx context.Context, cmd *CreateCommand) (string, error) {
	// check workspace exist
	gotWorkspace, err := h.workspaceClient.GetWorkspace(ctx, &proto.GetWorkspaceRequest{
		Id: cmd.WorkspaceID,
	})
	if err != nil {
		// TODO how to check not found error? grpc status error define in internal
		/*if gerr, ok := err.(*status.Error); ok {
			if gerr.GRPCStatus().Code() == codes.NotFound {
				return "", errors.NewNotFoundError("workspace", cmd.WorkspaceID)
			}
		}*/
		return "", err
	}
	if gotWorkspace.Workspace == nil {
		return "", errors.NewInternalError(fmt.Errorf("get workspace response nil object"))
	}

	// add workspace data
	var volumes []domain.Volume
	if gotWorkspace.Workspace.Storage != nil {
		if gotWorkspace.Workspace.Storage.Nfs != nil {
			volumes = append(volumes, domain.Volume{
				Type:              domain.VolumeTypeNFS,
				Name:              "bioos-workspace",
				Source:            gotWorkspace.Workspace.Storage.Nfs.MountPath,
				MountRelativePath: notebook.MountRelativePathWorkspaceData,
			})
		} else {
			log.Warnf("workspace %s has none storage", cmd.WorkspaceID)
		}
	} else {
		log.Warnf("workspace %s has none storage", cmd.WorkspaceID)
	}

	// add notebook ipynb
	if h.storageOpts.FileSystem != nil {
		volumes = append(volumes, domain.Volume{
			Type:              domain.VolumeTypeNFS,
			Name:              "bioos-notebook",
			Source:            path.Join(h.storageOpts.FileSystem.NotebookRootPath(), cmd.WorkspaceID), // TODO how to find out notebook store location?
			MountRelativePath: notebook.MountRelativePathNotebook,
		})
	}

	param := domain.CreateParam{
		WorkspaceID:  cmd.WorkspaceID,
		Image:        cmd.Image,
		ResourceSize: cmd.ResourceSize,
		Volumes:      volumes,
	}
	do, err := h.factory.New(&param)
	if err != nil {
		return "", err
	}
	if err = h.service.Create(ctx, do); err != nil {
		return "", err
	}
	return do.ID, nil
}
