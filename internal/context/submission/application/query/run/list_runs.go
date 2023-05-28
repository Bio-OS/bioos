package run

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListRunsQuery struct {
	WorkspaceID  string
	SubmissionID string
	Pg           *utils.Pagination
	Filter       *ListRunsFilter
}

type ListRunsHandler interface {
	Handle(context.Context, *ListRunsQuery) ([]*RunItem, int, error)
}

type listRunsHandler struct {
	runReadModel        ReadModel
	submissionReadModel submission.ReadModel
	workspaceClient     grpc.WorkspaceClient
}

func NewListRunsHandler(grpcFactory grpc.Factory, runReadModel ReadModel, submissionReadModel submission.ReadModel) ListRunsHandler {
	workspaceClient, err := grpcFactory.WorkspaceClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &listRunsHandler{
		runReadModel:        runReadModel,
		submissionReadModel: submissionReadModel,
		workspaceClient:     workspaceClient,
	}
}

func (l *listRunsHandler) Handle(ctx context.Context, query *ListRunsQuery) ([]*RunItem, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}

	if _, err := l.workspaceClient.GetWorkspace(ctx, &workspaceproto.GetWorkspaceRequest{Id: query.WorkspaceID}); err != nil {
		return nil, 0, apperrors.NewInternalError(err)
	}

	if err := submission.CheckSubmissionExist(ctx, l.submissionReadModel, query.WorkspaceID, query.SubmissionID); err != nil {
		return nil, 0, err
	}
	runs, err := l.runReadModel.ListRuns(ctx, query.SubmissionID, query.Pg, query.Filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := l.runReadModel.CountRuns(ctx, query.SubmissionID, query.Filter)
	if err != nil {
		return nil, 0, err
	}
	for _, run := range runs {
		statusCount, err := l.runReadModel.CountTasksResult(ctx, run.ID)
		if err != nil {
			return nil, 0, err
		}
		TaskStatus := Status{Count: 0}
		for _, v := range statusCount {
			TaskStatus.Count += v.Count
			switch v.Status {
			case consts.TaskSucceeded:
				TaskStatus.Succeeded += v.Count
			case consts.TaskRunning:
				TaskStatus.Running += v.Count
			case consts.TaskFailed:
				TaskStatus.Failed += v.Count
			case consts.TaskQueued:
				TaskStatus.Queued += v.Count
			case consts.TaskInitializing:
				TaskStatus.Initializing += v.Count
			case consts.TaskCancelled:
				TaskStatus.Cancelled += v.Count
			}
		}
		run.TaskStatus = TaskStatus
	}
	return runs, count, nil
}
