package submission

import "time"

// Submission ...
type Submission struct {
	ID                string
	Name              string
	Description       *string
	WorkflowID        string
	WorkflowVersionID string
	WorkspaceID       string
	DataModelID       *string
	DataModelRowIDs   []string
	Type              string
	Inputs            map[string]interface{}
	Outputs           map[string]interface{}
	ExposedOptions    string
	Status            string
	StartTime         time.Time
	FinishTime        *time.Time
}
