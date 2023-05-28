package convert

import submissionproto "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc/proto"

type ListRunsRequest struct {
	WorkspaceID  string   `path:"workspace_id"`
	SubmissionID string   `path:"submission_id"`
	Page         int      `query:"page"`
	Size         int      `query:"size"`
	OrderBy      string   `query:"orderBy"`
	SearchWord   string   `query:"searchWord"`
	Status       []string `query:"status,omitempty"`
	IDs          []string `query:"ids,omitempty"`
}

func (req *ListRunsRequest) ToGRPC() *submissionproto.ListRunsRequest {
	return &submissionproto.ListRunsRequest{
		WorkspaceID:  req.WorkspaceID,
		SubmissionID: req.SubmissionID,
		Page:         int32(req.Page),
		Size:         int32(req.Size),
		OrderBy:      req.OrderBy,
		SearchWord:   req.SearchWord,
		Status:       req.Status,
		Ids:          req.IDs,
	}
}

type ListRunsResponse struct {
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
	Items []RunItem `json:"items"`
}

func (resp *ListRunsResponse) FromGRPC(protoResp *submissionproto.ListRunsResponse) {
	resp.Page = int(protoResp.GetPage())
	resp.Size = int(protoResp.GetSize())
	resp.Total = int(protoResp.GetTotal())
	resp.Items = make([]RunItem, len(protoResp.GetItems()))
	for i, item := range protoResp.GetItems() {
		resp.Items[i] = RunItem{
			ID:          item.GetId(),
			Name:        item.GetName(),
			Status:      item.GetStatus(),
			StartTime:   item.GetStartTime(),
			FinishTime:  &item.FinishTime,
			Duration:    item.GetDuration(),
			EngineRunID: item.GetEngineRunID(),
			Inputs:      item.GetInputs(),
			Outputs:     item.GetOutputs(),
			TaskStatus: Status{
				Count:        item.GetTaskStatus().GetCount(),
				Pending:      item.GetTaskStatus().GetPending(),
				Succeeded:    item.GetTaskStatus().GetSucceeded(),
				Failed:       item.GetTaskStatus().GetFailed(),
				Running:      item.GetTaskStatus().GetRunning(),
				Cancelling:   item.GetTaskStatus().GetCancelling(),
				Cancelled:    item.GetTaskStatus().GetCancelled(),
				Queued:       item.GetTaskStatus().GetQueued(),
				Initializing: item.GetTaskStatus().GetInitializing(),
			},
			Log:     &item.Log,
			Message: &item.Message,
		}
	}

}

type RunItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	StartTime   int64   `json:"startTime"`
	FinishTime  *int64  `json:"finishTime"`
	Duration    int64   `json:"duration"`
	EngineRunID string  `json:"engineRunID"`
	Inputs      string  `json:"inputs"`
	Outputs     string  `json:"outputs"`
	TaskStatus  Status  `json:"taskStatus"`
	Log         *string `json:"log"`
	Message     *string `json:"message"`
}

type CancelRunRequest struct {
	WorkspaceID  string `path:"workspace_id"`
	SubmissionID string `path:"submission_id"`
	ID           string `path:"id"`
}

func (req *CancelRunRequest) ToGRPC() *submissionproto.CancelRunRequest {
	return &submissionproto.CancelRunRequest{
		WorkspaceID:  req.WorkspaceID,
		SubmissionID: req.SubmissionID,
		Id:           req.ID,
	}
}

type CancelRunResponse struct {
}

func (resp *CancelRunResponse) FromGRPC(protoResp *submissionproto.CancelRunResponse) {
	return
}

type ListTasksRequest struct {
	WorkspaceID  string `path:"workspace_id"`
	SubmissionID string `path:"submission_id"`
	RunID        string `path:"run_id"`
	Page         int    `query:"page"`
	Size         int    `query:"size"`
	OrderBy      string `query:"orderBy"`
}

func (req *ListTasksRequest) ToGRPC() *submissionproto.ListTasksRequest {
	return &submissionproto.ListTasksRequest{
		WorkspaceID:  req.WorkspaceID,
		SubmissionID: req.SubmissionID,
		RunID:        req.RunID,
		Page:         int32(req.Page),
		Size:         int32(req.Size),
		OrderBy:      req.OrderBy,
	}
}

type ListTasksResponse struct {
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Total int        `json:"total"`
	Items []TaskItem `json:"items"`
}

func (resp *ListTasksResponse) FromGRPC(protoResp *submissionproto.ListTasksResponse) {
	resp.Page = int(protoResp.GetPage())
	resp.Size = int(protoResp.GetSize())
	resp.Total = int(protoResp.GetTotal())
	resp.Items = make([]TaskItem, len(protoResp.GetItems()))
	for i, item := range protoResp.GetItems() {
		resp.Items[i] = TaskItem{
			Name:       item.GetName(),
			RunID:      item.GetRunID(),
			Status:     item.GetStatus(),
			StartTime:  item.GetStartTime(),
			FinishTime: &item.FinishTime,
			Duration:   item.GetDuration(),
			Stdout:     item.GetStdout(),
			Stderr:     item.GetStderr(),
		}
	}
}

type TaskItem struct {
	Name       string `json:"name"`
	RunID      string `json:"runID"`
	Status     string `json:"status"`
	StartTime  int64  `json:"startTime"`
	FinishTime *int64 `json:"finishTime"`
	Duration   int64  `json:"duration"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}
