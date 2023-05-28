package workspace

import (
	"time"

	"github.com/Bio-OS/bioos/pkg/utils"
)

type GetWorkspaceByIDQuery struct {
	ID string `validate:"required"`
}

type ListWorkspacesQuery struct {
	Pg     utils.Pagination
	Filter *ListWorkspacesFilter
}

type WorkspaceItem struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Storage     WorkspaceStorage
}

type WorkspaceStorage struct {
	NFS *NFSWorkspaceStorage
}

type NFSWorkspaceStorage struct {
	MountPath string
}

type ListWorkspacesFilter struct {
	SearchWord string
	Exact      bool
	IDs        []string
}

// Field for order.
const (
	OrderByName       = "Name"
	OrderByCreateTime = "CreatedAt"
)

type Queries struct {
	GetWorkspaceByID GetWorkspaceByIDQueryHandler
	ListWorkspaces   ListWorkspacesHandler
}

func NewQueries(workspaceReadModel WorkspaceReadModel) *Queries {
	return &Queries{
		GetWorkspaceByID: NewGetWorkspaceByIDHandler(workspaceReadModel),
		ListWorkspaces:   NewListWorkspacesHandler(workspaceReadModel),
	}
}
