package run

import (
	"github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/run"
	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type CancelRunCommand struct {
	WorkspaceID  string `validate:"required"`
	SubmissionID string `validate:"required"`
	ID           string `validate:"required"`
}

type Commands struct {
	CancelRun CancelRunHandler
}

func NewCommands(grpcFactory grpc.Factory, runRepo run.Repository, eventBus eventbus.EventBus, submissionReadModel submission.ReadModel, wesClient wes.Client) *Commands {
	service := run.NewService(grpcFactory, runRepo, eventBus, wesClient)
	return &Commands{
		CancelRun: NewCancelRunHandler(service, eventBus, submissionReadModel),
	}
}
