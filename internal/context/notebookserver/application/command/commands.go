package command

import (
	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/storage"
)

type Commands struct {
	Create CreateHandler
	Update UpdateHandler
	Switch SwitchHandler
	Delete DeleteHandler
}

func NewCommands(repo domain.Repository, factory *domain.Factory, runtime domain.Runtime, workspaceService proto.WorkspaceServiceServer, storageOpts *storage.Options, bus eventbus.EventBus) *Commands {
	svc := domain.NewService(repo, runtime)
	addEventHandle(bus, svc, factory, workspaceService, storageOpts)
	return &Commands{
		Create: NewCreateHandler(svc, factory, workspaceService, storageOpts),
		Update: NewUpdateHandler(svc),
		Switch: NewSwitchHandler(svc),
		Delete: NewDeleteHandler(svc),
	}
}
