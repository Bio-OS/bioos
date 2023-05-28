package query

import (
	"context"
	"net/url"
	"path"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetQuery struct {
	ID           string `validate:"required"`
	WorkspaceID  string `validate:"required"`
	EditNotebook string
}

type GetHandler interface {
	Handle(context.Context, *GetQuery) (*NotebookServer, error)
}

type getHandler struct {
	readModel ReadModel
	runtime   domain.Runtime
}

func NewGetHandler(readModel ReadModel, runtime domain.Runtime) GetHandler {
	return &getHandler{
		readModel: readModel,
		runtime:   runtime,
	}
}

func (r *getHandler) Handle(ctx context.Context, q *GetQuery) (*NotebookServer, error) {
	if err := validator.Validate(q); err != nil {
		return nil, err
	}
	settings, err := r.readModel.GetSettingsByID(ctx, q.WorkspaceID, q.ID)
	if err != nil {
		return nil, err
	}

	status, err := r.runtime.GetStatus(ctx, &domain.NotebookServer{
		ID:          q.ID,
		WorkspaceID: q.WorkspaceID,
	})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	res := &NotebookServer{
		Status:           status.Status,
		NotebookSettings: *settings,
	}
	if q.EditNotebook == "" {
		res.AccessURL = status.AccessURL
	} else { // locate to specified ipynb
		// TODO check notebook exist
		res.AccessURL = getAccessURL(status, q.EditNotebook)
	}
	return res, nil
}

func getAccessURL(status *domain.Status, editNotebook string) string {
	if status.Status == domain.ServerStatusRunning {
		u, _ := url.Parse(status.AccessURL)
		return u.JoinPath(getNotebookEditRelativePath(editNotebook)).String()
	} else if status.Status == domain.ServerStatusPending {
		u, _ := url.Parse(status.AccessURL)
		urlQuery := u.Query()
		if next := urlQuery.Get("next"); next != "" {
			next = path.Join(next, getNotebookEditRelativePath(editNotebook))
			urlQuery.Set("next", next)
			u.RawQuery = urlQuery.Encode()
		}
		return u.String()
	}
	return ""
}

func getNotebookEditRelativePath(editNotebook string) string {
	return path.Join(
		"lab",
		"tree",
		notebook.MountRelativePathNotebook,
		editNotebook+notebook.NotebookFileExt,
	)
}
