package query

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListQuery struct {
	WorkspaceID string `validate:"required"`
}

type ListHandler interface {
	Handle(context.Context, *ListQuery) ([]NotebookServer, error)
}

type listHandler struct {
	readModel ReadModel
	runtime   domain.Runtime
}

func NewListHandler(readModel ReadModel, runtime domain.Runtime) ListHandler {
	return &listHandler{
		readModel: readModel,
		runtime:   runtime,
	}
}

func (r *listHandler) Handle(ctx context.Context, q *ListQuery) ([]NotebookServer, error) {
	if err := validator.Validate(q); err != nil {
		return nil, err
	}
	settings, err := r.readModel.ListSettingsByWorkspace(ctx, q.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if len(settings) == 0 {
		return nil, nil
	}

	var res []NotebookServer
	for i := range settings {
		if settings[i] == nil {
			log.Warnf("workspace notebook server settings list %s has nil point", q.WorkspaceID)
			continue
		}
		status, err := r.runtime.GetStatus(ctx, &domain.NotebookServer{
			ID:          settings[i].ID,
			WorkspaceID: settings[i].WorkspaceID,
		})
		if err != nil {
			return nil, errors.NewInternalError(err)
		}
		res = append(res, NotebookServer{
			Status:           status.Status,
			NotebookSettings: *settings[i],
		})
	}
	return res, nil
}
