package workflow

const (
	WorkflowVersionPendingStatus = "Pending"
	WorkflowVersionSuccessStatus = "Success"
	WorkflowVersionFailedStatus  = "Failed"
)

const (
	WorkflowSourceGit  = "git"
	WorkflowSourceFile = "file"
)

const (
	WorkflowGitTag   = "gitTag"
	WorkflowGitURL   = "gitURL"
	WorkflowGitToken = "gitToken"
)

const (
	Language         = "WDL"
	VersionRegexpStr = "^version\\s+([\\w-._]+)"
)
