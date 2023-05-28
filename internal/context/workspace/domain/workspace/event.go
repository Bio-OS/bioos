package workspace

import (
	"encoding/json"
	"time"
)

const (
	WorkspaceCreated  string = "WorkspaceCreated"
	WorkspaceDeleted  string = "WorkspaceDeleted"
	WorkspaceImported string = "WorkspaceImported"

	ImportWorkspace string = "ImportWorkspace"
)

type WorkspaceEvent struct {
	WorkspaceID string
	Event       string
}

func NewWorkspaceCreatedEvent(workspaceID string) *WorkspaceEvent {
	return newWorkspaceEvent(workspaceID, WorkspaceCreated)
}

func NewWorkspaceDeletedEvent(workspaceID string) *WorkspaceEvent {
	return newWorkspaceEvent(workspaceID, WorkspaceDeleted)
}

func newWorkspaceEvent(workspaceID, event string) *WorkspaceEvent {
	return &WorkspaceEvent{
		WorkspaceID: workspaceID,
		Event:       event,
	}
}
func (e *WorkspaceEvent) EventType() string {
	return e.Event
}

func (e *WorkspaceEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *WorkspaceEvent) Delay() time.Duration {
	return 0
}

func NewWorkspaceEventFromPayload(data []byte) (*WorkspaceEvent, error) {
	ret := &WorkspaceEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type ImportWorkspaceEvent struct {
	WorkspaceID string
	FileName    string
	Storage     Storage
	Event       string
}

func NewImportWorkspaceEvent(workspaceID, fileName string, storage Storage) *ImportWorkspaceEvent {
	return &ImportWorkspaceEvent{
		WorkspaceID: workspaceID,
		FileName:    fileName,
		Storage:     storage,
		Event:       ImportWorkspace,
	}
}

func (e *ImportWorkspaceEvent) EventType() string {
	return e.Event
}

func (e *ImportWorkspaceEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *ImportWorkspaceEvent) Delay() time.Duration {
	return 0
}

func NewImportWorkspaceEventFromPayload(data []byte) (*ImportWorkspaceEvent, error) {
	ret := &ImportWorkspaceEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type WorkspaceImportedEvent struct {
	WorkspaceID   string
	ImportBaseDir string
}

func NewWorkspaceImportedEvent(workspaceID, baseDir string) *WorkspaceImportedEvent {
	return &WorkspaceImportedEvent{
		WorkspaceID:   workspaceID,
		ImportBaseDir: baseDir,
	}
}

func (e *WorkspaceImportedEvent) EventType() string {
	return WorkspaceImported
}

func (e *WorkspaceImportedEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *WorkspaceImportedEvent) Delay() time.Duration {
	return 0
}

func NewWorkspaceImportedEventFromPayload(data []byte) (*WorkspaceImportedEvent, error) {
	ret := &WorkspaceImportedEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
