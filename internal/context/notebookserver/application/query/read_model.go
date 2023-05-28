package query

import (
	"context"
)

type ReadModel interface {
	ListSettingsByWorkspace(context.Context, string) ([]*NotebookSettings, error)
	GetSettingsByID(ctx context.Context, workspaceID, id string) (*NotebookSettings, error)
}
