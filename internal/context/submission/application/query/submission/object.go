package submission

type SubmissionItem struct {
	ID                string
	Name              string
	Description       *string
	Type              string
	Status            string
	StartTime         int64
	FinishTime        *int64
	Duration          int64
	WorkflowID        string
	WorkflowVersionID string
	RunStatus         Status
	Entity            *Entity
	ExposedOptions    string
	InOutMaterial     *InOutMaterial
	WorkspaceID       string
}

type Entity struct {
	DataModelID     string
	DataModelRowIDs []string
	InputsTemplate  string
	OutputsTemplate string
}

type InOutMaterial struct {
	InputsMaterial  string
	OutputsMaterial string
}

type WorkflowVersion struct {
	ID        string
	VersionID string
}

type Status struct {
	Count      int64
	Pending    int64
	Succeeded  int64
	Failed     int64
	Running    int64
	Cancelling int64
	Cancelled  int64
}

type ListSubmissionsFilter struct {
	SearchWord string
	Exact      bool
	WorkflowID string
	Name       string
	Status     []string
	IDs        []string
}

// Field for order.
const (
	OrderByName      = "Name"
	OrderByStartTime = "StartTime"
)
