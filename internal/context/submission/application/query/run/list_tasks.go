package run

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/errors"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListTasksQuery struct {
	WorkspaceID  string
	SubmissionID string
	RunID        string
	Pg           *utils.Pagination
}

type ListTasksHandler interface {
	Handle(context.Context, *ListTasksQuery) ([]*TaskItem, int, error)
}

type listTasksHandler struct {
	runReadModel        ReadModel
	submissionReadModel submission.ReadModel
	workspaceClient     grpc.WorkspaceClient
}

func NewListTasksHandler(grpcFactory grpc.Factory, runReadModel ReadModel, submissionReadModel submission.ReadModel) ListTasksHandler {
	workspaceClient, err := grpcFactory.WorkspaceClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &listTasksHandler{
		runReadModel:        runReadModel,
		submissionReadModel: submissionReadModel,
		workspaceClient:     workspaceClient,
	}
}

func (l *listTasksHandler) Handle(ctx context.Context, query *ListTasksQuery) ([]*TaskItem, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}

	if _, err := l.workspaceClient.GetWorkspace(ctx, &workspaceproto.GetWorkspaceRequest{Id: query.WorkspaceID}); err != nil {
		return nil, 0, apperrors.NewInternalError(err)
	}

	if err := submission.CheckSubmissionExist(ctx, l.submissionReadModel, query.WorkspaceID, query.SubmissionID); err != nil {
		return nil, 0, err
	}
	runCount, err := l.runReadModel.CountRuns(ctx, query.SubmissionID, &ListRunsFilter{IDs: []string{query.RunID}})
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}
	if runCount != 1 {
		return nil, 0, fmt.Errorf("run id is not correct")
	}
	runs, err := l.runReadModel.ListTasks(ctx, query.RunID, query.Pg)
	if err != nil {
		return nil, 0, err
	}
	count, err := l.runReadModel.CountTasks(ctx, query.RunID)
	if err != nil {
		return nil, 0, err
	}
	return runs, count, nil
}
