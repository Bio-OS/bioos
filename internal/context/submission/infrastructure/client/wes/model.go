package wes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ListRunsRequest ...
type ListRunsRequest struct {
	PageSize     *int64
	PageToken    *string
	TagFilter    *string
	WorkflowType string `json:"workflow_type"`
}

// ListRunsResponse ...
type ListRunsResponse struct {
	Runs          []RunStatus `json:"runs"`
	NextPageToken string      `json:"next_page_token"`
}

// RunWorkflowRequest ...
type RunWorkflowRequest struct {
	RunRequest
	WorkflowAttachment map[string]string // the string is file path
}

// RunWorkflowResponse ...
type RunWorkflowResponse struct {
	RunID string `json:"run_id"`
}

// GetRunLogRequest ...
type GetRunLogRequest struct {
	RunID        string
	WorkflowType string `json:"workflow_type"`
}

// GetRunLogResponse ...
type GetRunLogResponse struct {
	RunID    string                 `json:"run_id"`
	Request  RunRequest             `json:"request"`
	State    RunState               `json:"state"`
	RunLog   Log                    `json:"run_log"`
	TaskLogs []Log                  `json:"task_logs"`
	Outputs  map[string]interface{} `json:"outputs"`
}

// CancelRunRequest ...
type CancelRunRequest struct {
	RunID        string
	WorkflowType string `json:"workflow_type"`
}

// CancelRunResponse ...
type CancelRunResponse struct {
	RunID string `json:"run_id"`
}

// RunRequest ...
type RunRequest struct {
	WorkflowParams           map[string]interface{} `json:"workflow_params"`
	WorkflowType             string                 `json:"workflow_type"`
	WorkflowTypeVersion      string                 `json:"workflow_type_version"`
	Tags                     map[string]interface{} `json:"tags"`
	WorkflowEngineParameters map[string]interface{} `json:"workflow_engine_parameters"`
}

// RunStatus ...
type RunStatus struct {
	RunID string   `json:"run_id"`
	State RunState `json:"state"`
}

// RunState ...
type RunState string

// run state enum
const (
	RunStateUnknown       RunState = "UNKNOWN"
	RunStateQueued        RunState = "QUEUED"
	RunStateInitializing  RunState = "INITIALIZING"
	RunStateRunning       RunState = "RUNNING"
	RunStatePaused        RunState = "PAUSED"
	RunStateComplete      RunState = "COMPLETE"
	RunStateExecutorError RunState = "EXECUTOR_ERROR"
	RunStateSystemError   RunState = "SYSTEM_ERROR"
	RunStateCanceled      RunState = "CANCELED"
	RunStateCanceling     RunState = "CANCELING"
)

// Log ...
type Log struct {
	Name      string   `json:"name"`
	Cmd       []string `json:"cmd"`
	StartTime *Time    `json:"start_time"`
	EndTime   *Time    `json:"end_time"`
	Stdout    string   `json:"stdout"`
	Stderr    string   `json:"stderr"`
	Log       string   `json:"log"`
	ExitCode  *int32   `json:"exit_code"`
}

// ErrorResp ...
type ErrorResp struct {
	Msg        string `json:"msg"`
	StatusCode int32  `json:"status_code"`
}

// Error ...
func (e ErrorResp) Error() string {
	return fmt.Sprintf("Msg: %s, StatusCode: %d", e.Msg, e.StatusCode)
}

// IsNotFound ...
func IsNotFound(err error) bool {
	var wesErr ErrorResp
	return errors.As(err, &wesErr) && wesErr.StatusCode == http.StatusNotFound
}

// IsBadRequest ...
func IsBadRequest(err error) bool {
	var wesErr ErrorResp
	return errors.As(err, &wesErr) && wesErr.StatusCode == http.StatusBadRequest
}

func newBadRequestError(msg string) error {
	return ErrorResp{
		Msg:        msg,
		StatusCode: http.StatusBadRequest,
	}
}

type listRunsResponseWithError struct {
	ListRunsResponse `json:",omitempty,inline"`
	ErrorResp        `json:",omitempty,inline"`
}

type runWorkflowResponseWithError struct {
	RunWorkflowResponse `json:",omitempty,inline"`
	ErrorResp           `json:",omitempty,inline"`
}

type getRunLogResponseWithError struct {
	GetRunLogResponse `json:",omitempty,inline"`
	ErrorResp         `json:",omitempty,inline"`
}

type cancelRunResponseWithError struct {
	CancelRunResponse `json:",omitempty,inline"`
	ErrorResp         `json:",omitempty,inline"`
}

// Time ...
type Time time.Time

// UnmarshalJSON ...
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	asTime, err := time.Parse(time.RFC3339Nano, strings.Trim(string(data), "\""))
	*t = Time(asTime)
	return
}

// MarshalJSON ...
func (t *Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, time.RFC3339Nano)
	b = append(b, '"')
	return b, nil
}

// String ...
func (t *Time) String() string {
	return time.Time(*t).Format(time.RFC3339Nano)
}

// Time ...
func (t *Time) Time() time.Time {
	return time.Time(*t)
}

// PointTime ...
func (t *Time) PointTime() *time.Time {
	return (*time.Time)(t)
}
