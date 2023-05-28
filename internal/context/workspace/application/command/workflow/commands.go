package workflow

import (
	workflowquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

type Commands struct {
	Create CreateHandler
	Delete DeleteHandler
	Update UpdateHandler
}

func NewCommands(repo workflow.Repository, workflowReadModel workflowquery.ReadModel, factory *workflow.Factory, workspaceReadModel workspace.WorkspaceReadModel, bus eventbus.EventBus, womtoolPath string) *Commands {
	service := workflow.NewService(repo, workflowReadModel, bus, factory, womtoolPath)
	return &Commands{
		Create: NewCreateHandler(service, workflowReadModel, workspaceReadModel),
		Delete: NewDeleteHandler(service, workspaceReadModel),
		Update: NewUpdateHandler(service, workflowReadModel, workspaceReadModel),
	}
}
