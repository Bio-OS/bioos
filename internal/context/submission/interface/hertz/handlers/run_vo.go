package handlers

type ListRunsRequest struct {
	WorkspaceID  string   `path:"workspace_id"`
	SubmissionID string   `path:"submission_id"`
	Page         int      `query:"page"`
	Size         int      `query:"size"`
	OrderBy      string   `query:"orderBy"`
	SearchWord   string   `query:"searchWord"`
	Status       []string `query:"status"`
	IDs          []string `query:"ids"`
}

type ListRunsResponse struct {
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
	Items []RunItem `json:"items"`
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

type ListTasksRequest struct {
	WorkspaceID  string `path:"workspace_id"`
	SubmissionID string `path:"submission_id"`
	RunID        string `path:"run_id"`
	Page         int    `query:"page"`
	Size         int    `query:"size"`
	OrderBy      string `query:"orderBy"`
}

type ListTasksResponse struct {
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Total int        `json:"total"`
	Items []TaskItem `json:"items"`
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
