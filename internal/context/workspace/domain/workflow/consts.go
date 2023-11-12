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
	VersionWDLRegexpStr      = "^version\\s+([\\w-._]+)"
	VersionNextflowRegexpStr = "^nextflow.enable.dsl\\s+=\\s+(\\d)"
)

const (
	LanguageWDL      Language = "WDL"
	LanguageNextflow Language = "Nextflow"
)
