package workspace

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/pkg/errors"
)

func CheckWorkspaceExist(ctx context.Context, readModel WorkspaceReadModel, workspaceID string) *errors.AppError {
	filter := ListWorkspacesFilter{
		IDs: []string{workspaceID},
	}
	wsCount, err := readModel.CountWorkspaces(ctx, &filter)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("find workspace fail: %w", err))
	}
	if wsCount < 1 {
		return errors.NewNotFoundError("workspace", workspaceID)
	}
	return nil
}
