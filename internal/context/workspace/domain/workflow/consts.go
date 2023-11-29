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
	VersionRegexpStrWDL      = "^version\\s+([\\w-._]+)"
	VersionRegexpStrNextflow = "^nextflow.enable.dsl\\s+=\\s+(\\d)"
)

const (
	LanguageWDL       = "WDL"
	LanguageNextflow  = "Nextflow"
	LanguageCWL       = "CWL"
	LanguageSnakemake = "Snakemake"
)
