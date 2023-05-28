package query

import "github.com/Bio-OS/bioos/internal/context/notebookserver/domain"

type Queries struct {
	List ListHandler
	Get  GetHandler
}

func NewQueries(readModel ReadModel, runtime domain.Runtime) *Queries {
	return &Queries{
		List: NewListHandler(readModel, runtime),
		Get:  NewGetHandler(readModel, runtime),
	}
}
