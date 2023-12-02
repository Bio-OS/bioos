package run

type RunItem struct {
	ID           string
	Name         string
	Status       string
	StartTime    int64
	FinishTime   *int64
	Duration     int64
	EngineRunID  string
	Inputs       string
	Outputs      string
	TaskStatus   Status
	Log          *string
	Message      *string
	WorkflowType string
}

type TaskItem struct {
	Name       string
	RunID      string
	Status     string
	StartTime  int64
	FinishTime *int64
	Duration   int64
	Stdout     string
	Stderr     string
}

type Status struct {
	Count        int64
	Succeeded    int64
	Failed       int64
	Running      int64
	Cancelling   int64
	Cancelled    int64
	Queued       int64
	Initializing int64
	Pending      int64
}

type ListRunsFilter struct {
	SearchWord string
	Exact      bool
	Status     []string
	IDs        []string
}

// StatusCount ...
type StatusCount struct {
	Count  int64
	Status string
}

// Field for order.
const (
	OrderByName      = "Name"
	OrderByStartTime = "StartTime"
)
