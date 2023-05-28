package run

import (
	"github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type Queries struct {
	ListRuns        ListRunsHandler
	ListTasks       ListTasksHandler
	CountRunsResult CountRunsResultHandler
}

func NewQueries(grpcFactory grpc.Factory, runReadModel ReadModel, submissionReadModel submission.ReadModel) *Queries {
	return &Queries{
		ListRuns:        NewListRunsHandler(grpcFactory, runReadModel, submissionReadModel),
		ListTasks:       NewListTasksHandler(grpcFactory, runReadModel, submissionReadModel),
		CountRunsResult: NewCountRunsResultHandler(runReadModel),
	}
}
