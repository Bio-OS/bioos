package submission

import (
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type Queries struct {
	List  ListHandler
	Check CheckHandler
}

func NewQueries(grpcFactory grpc.Factory, submissionReadModel ReadModel) *Queries {
	return &Queries{
		List:  NewListHandler(grpcFactory, submissionReadModel),
		Check: NewCheckHandler(grpcFactory, submissionReadModel),
	}
}
