package notebook

import "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"

type Queries struct {
	List ListHandler
	Get  GetHandler
}

func NewQueries(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) *Queries {
	return &Queries{
		List: NewListHandler(readModel, workspaceReadModel),
		Get:  NewGetHandler(readModel, workspaceReadModel),
	}
}
