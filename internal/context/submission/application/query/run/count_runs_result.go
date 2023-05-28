package run

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CountRunsResultQuery struct {
	SubmissionID string
}

type CountRunsResultHandler interface {
	Handle(context.Context, *CountRunsResultQuery) (*Status, error)
}

type countRunsResultHandler struct {
	readModel ReadModel
}

func NewCountRunsResultHandler(readModel ReadModel) CountRunsResultHandler {
	return &countRunsResultHandler{
		readModel: readModel,
	}
}

func (l *countRunsResultHandler) Handle(ctx context.Context, query *CountRunsResultQuery) (*Status, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}
	statusCount, err := l.readModel.CountRunsResult(ctx, query.SubmissionID)
	if err != nil {
		return nil, err
	}
	RunStatus := &Status{Count: 0}
	for _, v := range statusCount {
		RunStatus.Count += v.Count
		switch v.Status {
		case consts.RunSucceeded:
			RunStatus.Succeeded += v.Count
		case consts.RunRunning:
			RunStatus.Running += v.Count
		case consts.RunFailed:
			RunStatus.Failed += v.Count
		case consts.RunPending:
			RunStatus.Pending += v.Count
		case consts.RunCancelling:
			RunStatus.Cancelling += v.Count
		case consts.RunCancelled:
			RunStatus.Cancelled += v.Count
		}
	}
	return RunStatus, nil
}
