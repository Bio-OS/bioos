package submission

type SubmissionItem struct {
	ID                string
	Name              string
	Description       *string
	Type              string
	Status            string
	Language          string
	StartTime         int64
	FinishTime        *int64
	Duration          int64
	WorkflowID        string
	WorkflowVersionID string
	RunStatus         Status
	Entity            *Entity
	ExposedOptions    ExposedOptions
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

type ExposedOptions struct {
	ReadFromCache bool
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
	Language   []string
}

// Field for order.
const (
	OrderByName      = "Name"
	OrderByStartTime = "StartTime"
)
