package workflow

import (
	"encoding/json"
	"time"

	"github.com/Bio-OS/bioos/pkg/schema"
)

const (
	WorkflowCreated = "WorkflowCreated"
	WorkflowDeleted = "WorkflowDeleted"

	WorkflowVersionAdded = "WorkflowVersionAdded"

	ImportWorkflows = "ImportWorkflows"
)

type WorkflowEvent struct {
	WorkspaceID string
	WorkflowID  string
	Event       string
}

func NewWorkflowCreatedEvent(workspaceID, workflowID string) *WorkflowEvent {
	return newWorkflowEvent(workspaceID, workflowID, WorkflowCreated)
}

func NewWorkflowDeletedEvent(workspaceID, workflowID string) *WorkflowEvent {
	return newWorkflowEvent(workspaceID, workflowID, WorkflowDeleted)
}

func newWorkflowEvent(workspaceID, workflowID, event string) *WorkflowEvent {
	return &WorkflowEvent{
		WorkspaceID: workspaceID,
		WorkflowID:  workflowID,
		Event:       event,
	}
}

func (e *WorkflowEvent) EventType() string {
	return e.Event
}

func (e *WorkflowEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *WorkflowEvent) Delay() time.Duration {
	return 0
}

func NewWorkflowEventFromPayload(data []byte) (*WorkflowEvent, error) {
	ret := &WorkflowEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type WorkflowVersionAddedEvent struct {
	WorkspaceID       string
	WorkflowID        string
	WorkflowVersionID string
	GitRepo           string
	GitTag            string
	GitToken          string
	//used to import workflow version from file. its priority is over Git
	FilesBaseDir string
}

func NewWorkflowVersionAddedEvent(workspaceID, workflowID, versionID, repo, tag, token, filesBaseDir string) *WorkflowVersionAddedEvent {
	return &WorkflowVersionAddedEvent{
		WorkspaceID:       workspaceID,
		WorkflowID:        workflowID,
		WorkflowVersionID: versionID,
		GitRepo:           repo,
		GitTag:            tag,
		GitToken:          token,
		FilesBaseDir:      filesBaseDir,
	}
}

func (e *WorkflowVersionAddedEvent) EventType() string {
	return WorkflowVersionAdded
}

func (e *WorkflowVersionAddedEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *WorkflowVersionAddedEvent) Delay() time.Duration {
	return 0
}

func NewWorkflowVersionAddedEventFromPayload(data []byte) (*WorkflowVersionAddedEvent, error) {
	ret := &WorkflowVersionAddedEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type ImportWorkflowsEvent struct {
	WorkspaceID       string
	Schemas           []schema.WorkflowTypedSchema
	ImportFileBaseDir string
}

func NewImportWorkflowsEvent(workspaceID, baseDir string, schemas []schema.WorkflowTypedSchema) *ImportWorkflowsEvent {
	return &ImportWorkflowsEvent{
		WorkspaceID:       workspaceID,
		Schemas:           schemas,
		ImportFileBaseDir: baseDir,
	}
}

func (e *ImportWorkflowsEvent) EventType() string {
	return ImportWorkflows
}

func (e *ImportWorkflowsEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *ImportWorkflowsEvent) Delay() time.Duration {
	return 0
}

func NewImportWorkflowsEventFromPayload(data []byte) (*ImportWorkflowsEvent, error) {
	ret := &ImportWorkflowsEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
