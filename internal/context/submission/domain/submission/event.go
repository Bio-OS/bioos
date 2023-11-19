package submission

import (
	"encoding/json"
	"github.com/Bio-OS/bioos/pkg/utils"
	"time"
)

const (
	CreateSubmission        = "CreateSubmission"
	CancelSubmission        = "CancelSubmission"
	DeleteSubmission        = "DeleteSubmission"
	CascadeDeleteSubmission = "CascadeDeleteSubmission"
	SyncSubmission          = "SyncSubmission"

	CreateRuns = "CreateRuns"
	SubmitRun  = "SubmitRun"
	SyncRun    = "SyncRun"
	CancelRun  = "CancelRun"
	DeleteRun  = "DeleteRun"
)

type EventSubmission struct {
	SubmissionID  string
	Event         string
	DelayDuration time.Duration
}

func NewCancelSubmissionEvent(submissionID string, duration time.Duration) *EventSubmission {
	return &EventSubmission{
		SubmissionID:  submissionID,
		Event:         CancelSubmission,
		DelayDuration: duration,
	}
}

func NewDeleteSubmissionEvent(submissionID string, duration time.Duration) *EventSubmission {
	return &EventSubmission{
		SubmissionID:  submissionID,
		Event:         DeleteSubmission,
		DelayDuration: duration,
	}
}

func NewSyncSubmissionEvent(submissionID string) *EventSubmission {
	return &EventSubmission{
		SubmissionID: submissionID,
		Event:        SyncSubmission,
	}
}

func (e *EventSubmission) EventType() string {
	return e.Event
}

func (e *EventSubmission) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *EventSubmission) Delay() time.Duration {
	return e.DelayDuration
}

type CreateEvent struct {
	WorkspaceID             string
	SubmissionID            string
	SourceWorkflowID        string
	SourceWorkflowVersionID string
	SourceDataModelID       *string
	SourceDataModelRowIDs   []string
	Language                string
}

func NewCreateEvent(workspaceID, submissionID, workflowID, workflowVersionID string, language string, sourceDataModelID *string, sourceDataModelRowIDs []string) *CreateEvent {
	return &CreateEvent{
		WorkspaceID:             workspaceID,
		SubmissionID:            submissionID,
		SourceWorkflowID:        workflowID,
		SourceWorkflowVersionID: workflowVersionID,
		SourceDataModelID:       sourceDataModelID,
		SourceDataModelRowIDs:   sourceDataModelRowIDs,
		Language:                language,
	}
}

func (e *CreateEvent) EventType() string {
	return CreateSubmission
}

func (e *CreateEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *CreateEvent) Delay() time.Duration {
	return 0
}

func NewEventFromPayload(data []byte) (*EventSubmission, error) {
	res := &EventSubmission{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

func NewCreateEventFromPayload(data []byte) (*CreateEvent, error) {
	res := &CreateEvent{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

type CascadeDeleteSubmissionEvent struct {
	WorkspaceID string
	Workflow    *string
}

func NewEventCascadeDeleteSubmission(workspaceID string, workflow *string) *CascadeDeleteSubmissionEvent {
	return &CascadeDeleteSubmissionEvent{
		WorkspaceID: workspaceID,
		Workflow:    workflow,
	}
}

func (e *CascadeDeleteSubmissionEvent) EventType() string {
	return CascadeDeleteSubmission
}

func (e *CascadeDeleteSubmissionEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *CascadeDeleteSubmissionEvent) Delay() time.Duration {
	return 0
}

func NewEventCascadeDeleteSubmissionFromPayload(data []byte) (*CascadeDeleteSubmissionEvent, error) {
	res := &CascadeDeleteSubmissionEvent{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

type EventCreateRuns struct {
	WorkspaceID     string
	SubmissionID    string
	InputsTemplate  map[string]interface{}
	OutputsTemplate map[string]interface{}
	SubmissionType  string // filePath or dataModel
	DataModelID     *string
	DataModelRowIDs []string

	RunConfig *RunConfig
}

func (e *EventCreateRuns) InAndOutTemplate() map[string]interface{} {
	return utils.MergeMapInterface(e.InputsTemplate, e.OutputsTemplate)
}

type RunConfig struct {
	Language                 string
	WorkflowContents         map[string]string
	MainWorkflowFilePath     string
	WorkflowEngineParameters map[string]interface{}
	Version                  string
	WorkflowURL              string
}

func NewEventCreateRuns(workspaceID, submissionID, submissionType string, inputs, outputs map[string]interface{}, dataModelID *string, DataModelRowIDs []string, runConfig *RunConfig) *EventCreateRuns {
	return &EventCreateRuns{
		WorkspaceID:     workspaceID,
		SubmissionID:    submissionID,
		SubmissionType:  submissionType,
		InputsTemplate:  inputs,
		OutputsTemplate: outputs,
		DataModelID:     dataModelID,
		DataModelRowIDs: DataModelRowIDs,
		RunConfig:       runConfig,
	}
}

func (e *EventCreateRuns) EventType() string {
	return CreateRuns
}

func (e *EventCreateRuns) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *EventCreateRuns) Delay() time.Duration {
	return 0
}

func NewEventCreateRunFromPayload(data []byte) (*EventCreateRuns, error) {
	res := &EventCreateRuns{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

type EventSubmitRun struct {
	RunID     string
	RunConfig *RunConfig
}

func (e *EventSubmitRun) WorkflowType() string {
	return e.RunConfig.Language
}

func NewEventSubmitRun(runID string, runConfig *RunConfig) *EventSubmitRun {
	return &EventSubmitRun{
		RunID:     runID,
		RunConfig: runConfig,
	}
}

func (e *EventSubmitRun) EventType() string {
	return SubmitRun
}

func (e *EventSubmitRun) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *EventSubmitRun) Delay() time.Duration {
	return 0
}

func NewEventSubmitRunFromPayload(data []byte) (*EventSubmitRun, error) {
	res := &EventSubmitRun{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

type EventRun struct {
	RunID         string
	EventTyp      string
	DelayDuration time.Duration
	WorkflowType  string
}

func (e *EventRun) EventType() string {
	return e.EventTyp
}

func (e *EventRun) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *EventRun) Delay() time.Duration {
	return e.DelayDuration
}

func NewEventRunFromPayload(data []byte) (*EventRun, error) {
	res := &EventRun{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	return res, nil
}

func NewEventDeleteRun(runID string) *EventRun {
	return &EventRun{
		RunID:         runID,
		EventTyp:      DeleteRun,
		DelayDuration: 0,
	}
}

func NewEventSyncRun(runID string, workflowType string, delayDuration time.Duration) *EventRun {
	return &EventRun{
		RunID:         runID,
		EventTyp:      SyncRun,
		DelayDuration: delayDuration,
		WorkflowType:  workflowType,
	}
}

func NewEventCancelRun(runID string) *EventRun {
	return &EventRun{
		RunID:    runID,
		EventTyp: CancelRun,
	}
}
