package datamodel

import (
	"encoding/json"
	"time"

	"github.com/Bio-OS/bioos/pkg/schema"
)

const (
	ImportDataModels = "ImportDataModels"
)

type ImportDataModelsEvent struct {
	WorkspaceID       string
	Schemas           []schema.DataModelTypedSchema
	ImportFileBaseDir string
}

func NewImportDataModelsEvent(workspaceID, baseDir string, schemas []schema.DataModelTypedSchema) *ImportDataModelsEvent {
	return &ImportDataModelsEvent{
		WorkspaceID:       workspaceID,
		Schemas:           schemas,
		ImportFileBaseDir: baseDir,
	}
}

func (e *ImportDataModelsEvent) EventType() string {
	return ImportDataModels
}

func (e *ImportDataModelsEvent) Payload() []byte {
	payload, _ := json.Marshal(e)
	return payload
}

func (e *ImportDataModelsEvent) Delay() time.Duration {
	return 0
}

func NewImportDataModelsEventFromPayload(data []byte) (*ImportDataModelsEvent, error) {
	ret := &ImportDataModelsEvent{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
