package submission

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/pkg/errors"
)

func CheckSubmissionExist(ctx context.Context, readModel ReadModel, workspaceID, id string) *errors.AppError {
	filter := &ListSubmissionsFilter{
		IDs: []string{id},
	}
	wsCount, err := readModel.CountSubmissions(ctx, workspaceID, filter)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("find submission fail: %w", err))
	}
	if wsCount < 1 {
		return errors.NewNotFoundError("submission", id)
	}
	return nil
}
