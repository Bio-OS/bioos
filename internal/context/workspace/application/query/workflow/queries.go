package workflow

import (
	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
)

type Queries struct {
	GetByID       GetHandler
	GetFile       GetFileHandler
	GetVersion    GetVersionHandler
	ListWorkflows ListHandler
	ListFiles     ListFilesHandler
	ListVersions  ListVersionsHandler
}

func NewQueries(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) *Queries {
	return &Queries{
		GetByID:       NewGetHandler(readModel, workspaceReadModel),
		GetFile:       NewGetFileHandler(readModel, workspaceReadModel),
		GetVersion:    NewGetVersionHandler(readModel, workspaceReadModel),
		ListWorkflows: NewListHandler(readModel, workspaceReadModel),
		ListFiles:     NewListFilesHandler(readModel, workspaceReadModel),
		ListVersions:  NewListVersionsHandler(readModel, workspaceReadModel),
	}
}
