package workspace

import (
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

type CreateWorkspaceCommand struct {
	Name        string `validate:"required,resName"`
	Description string `validate:"required,workspaceDesc"`
	Storage     WorkspaceStorage
}

type DeleteWorkspaceCommand struct {
	ID string `validate:"required"`
}

type UpdateWorkspaceCommand struct {
	ID          string  `validate:"required"`
	Name        *string `validate:"omitempty,resName"`
	Description *string `validate:"omitempty,workspaceDesc"`
}

type WorkspaceStorage struct {
	NFS *NFSWorkspaceStorage
}

type NFSWorkspaceStorage struct {
	MountPath string `validate:"required,nfsMountPath"`
}

type ImportWorkspaceCommand struct {
	ID       string `validate:"required"`
	FileName string `validate:"required"`
	Storage  WorkspaceStorage
}

type Commands struct {
	CreateWorkspace CreateWorkspaceHandler
	ImportWorkspace ImportWorkspaceHandler
	DeleteWorkspace DeleteWorkspaceHandler
	UpdateWorkspace UpdateWorkspaceHandler
}

func NewCommands(workspaceRepo workspace.Repository, eventRepo eventbus.EventRepository, workspaceFactory *workspace.Factory, eventBus eventbus.EventBus) *Commands {
	service := workspace.NewService(workspaceRepo, eventRepo, eventBus, *workspaceFactory)
	addEventHandle(eventBus, workspaceRepo, eventRepo, workspaceFactory)
	return &Commands{
		CreateWorkspace: NewCreateWorkspaceHandler(workspaceRepo, workspaceFactory, eventBus),
		ImportWorkspace: NewImportWorkspaceHandler(service, workspaceFactory),
		DeleteWorkspace: NewDeleteWorkspaceHandler(workspaceRepo, eventBus),
		UpdateWorkspace: NewUpdateWorkspaceHandler(workspaceRepo, eventBus),
	}
}
