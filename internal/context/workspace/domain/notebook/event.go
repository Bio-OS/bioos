package notebook

import (
	"encoding/json"
	"time"

	"github.com/Bio-OS/bioos/pkg/schema"
)

const (
	ImportNotebooks = "ImportNotebooks"
)

type ImportNotebooksEvent struct {
	WorkspaceID       string
	Schema            schema.NotebookTypedSchema
	ImportFileBaseDir string
}

func NewImportNotebooksEvent(workspaceID, baseDir string, schema schema.NotebookTypedSchema) *ImportNotebooksEvent {
	return &ImportNotebooksEvent{
		WorkspaceID:       workspaceID,
		Schema:            schema,
		ImportFileBaseDir: baseDir,
	}
}

func (e *ImportNotebooksEvent) EventType() string {
	return ImportNotebooks
}

func (e *ImportNotebooksEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *ImportNotebooksEvent) Delay() time.Duration {
	return 0
}

func NewImportNotebooksEventFromPayload(data []byte) (*ImportNotebooksEvent, error) {
	ret := &ImportNotebooksEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
