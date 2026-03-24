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
	WorkflowLanguageWDL       = "WDL"
	WorkflowLanguageCWL       = "CWL"
	WorkflowLanguageSnakemake = "SMK"
	WorkflowLanguageNextflow  = "NFL"
)

const (
	Language         = "WDL"
	VersionRegexpStr = "^version\\s+([\\w-._]+)"
)
