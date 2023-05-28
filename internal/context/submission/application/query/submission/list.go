package submission

import (
	"context"

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListQuery struct {
	WorkspaceID string
	Pg          *utils.Pagination
	Filter      *ListSubmissionsFilter
}

type ListHandler interface {
	Handle(context.Context, *ListQuery) ([]*SubmissionItem, int, error)
}

type listHandler struct {
	submissionReadModel ReadModel
	workspaceClient     grpc.WorkspaceClient
}

func NewListHandler(grpcFactory grpc.Factory, submissionReadModel ReadModel) ListHandler {
	workspaceClient, err := grpcFactory.WorkspaceClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &listHandler{
		submissionReadModel: submissionReadModel,
		workspaceClient:     workspaceClient,
	}
}

func (l *listHandler) Handle(ctx context.Context, query *ListQuery) ([]*SubmissionItem, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}

	if _, err := l.workspaceClient.GetWorkspace(ctx, &workspaceproto.GetWorkspaceRequest{Id: query.WorkspaceID}); err != nil {
		return nil, 0, apperrors.NewInternalError(err)
	}

	subs, err := l.submissionReadModel.ListSubmissions(ctx, query.WorkspaceID, query.Pg, query.Filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := l.submissionReadModel.CountSubmissions(ctx, query.WorkspaceID, query.Filter)
	if err != nil {
		return nil, 0, err
	}
	return subs, count, nil
}
