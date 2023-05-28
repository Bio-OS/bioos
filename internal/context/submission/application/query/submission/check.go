package submission

import (
	"context"

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/errors"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CheckQuery struct {
	WorkspaceID string
	Name        string
}

type CheckHandler interface {
	Handle(context.Context, *CheckQuery) (bool, error)
}

type checkHandler struct {
	readModel       ReadModel
	workspaceClient grpc.WorkspaceClient
}

func NewCheckHandler(grpcFactory grpc.Factory, readModel ReadModel) CheckHandler {
	workspaceClient, err := grpcFactory.WorkspaceClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &checkHandler{
		readModel:       readModel,
		workspaceClient: workspaceClient,
	}
}

func (l *checkHandler) Handle(ctx context.Context, query *CheckQuery) (bool, error) {
	if err := validator.Validate(query); err != nil {
		return false, err
	}

	if _, err := l.workspaceClient.GetWorkspace(ctx, &workspaceproto.GetWorkspaceRequest{Id: query.WorkspaceID}); err != nil {
		return false, apperrors.NewInternalError(err)
	}

	count, err := l.readModel.CountSubmissions(ctx, query.WorkspaceID, &ListSubmissionsFilter{Name: query.Name})
	if err != nil {
		return false, errors.NewInternalError(err)
	}
	return count > 0, nil
}
