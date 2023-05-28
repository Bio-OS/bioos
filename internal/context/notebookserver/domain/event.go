package domain

import (
	"encoding/json"
	"time"

	"github.com/Bio-OS/bioos/pkg/schema"
)

const (
	ImportNotebookServers = "ImportNotebookServers"
)

type ImportNotebookServersEvent struct {
	WorkspaceID       string
	Schema            schema.NotebookTypedSchema
	ImportFileBaseDir string
}

func NewImportNotebookServersEvent(workspaceID, baseDir string, schema schema.NotebookTypedSchema) *ImportNotebookServersEvent {
	return &ImportNotebookServersEvent{
		WorkspaceID:       workspaceID,
		Schema:            schema,
		ImportFileBaseDir: baseDir,
	}
}

func (e *ImportNotebookServersEvent) EventType() string {
	return ImportNotebookServers
}

func (e *ImportNotebookServersEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *ImportNotebookServersEvent) Delay() time.Duration {
	return 0
}

func NewImportNotebookServersEventFromPayload(data []byte) (*ImportNotebookServersEvent, error) {
	ret := &ImportNotebookServersEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
