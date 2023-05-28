package notebook

import (
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

type Commands struct {
	Create CreateHandler
	Update UpdateHandler
	Delete DeleteHandler
}

func NewCommands(repo notebook.Repository, workspaceReadModel workspace.WorkspaceReadModel, eb eventbus.EventBus, readModel query.ReadModel, factory *notebook.Factory) *Commands {
	svc := notebook.NewService(repo)
	addEventHandle(eb, readModel, svc, factory)
	return &Commands{
		Create: NewCreateHandler(svc, workspaceReadModel),
		Update: NewUpdateHandler(svc, workspaceReadModel),
		Delete: NewDeleteHandler(svc, workspaceReadModel),
	}
}
