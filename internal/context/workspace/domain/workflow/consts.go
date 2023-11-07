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
	VersionRegexpStr = "^version\\s+([\\w-._]+)"
)

const (
	LanguageWDL      Language = "WDL"
	LanguageNextflow Language = "Nextflow"
)
