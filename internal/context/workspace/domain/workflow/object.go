package workflow

import (
	"time"

	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// Workflow workflow entity
type Workflow struct {
	// ID is the unique identifier of the workflow
	ID string
	// Name is the name of the workflow
	Name string
	// WorkspaceID is the ID of the workspace
	WorkspaceID string
	// Description is the description of the workflow
	Description string
	// LatestVersion is the latest version of the workflow
	LatestVersion string
	// Versions is the versions of the workflow
	Versions map[string]*WorkflowVersion
	// CreatedAt is the create time of workflow version
	CreatedAt time.Time
	// UpdatedAt is the update time of workflow version
	UpdatedAt time.Time
}

func (w *Workflow) UpdateName(name string) {
	if w.Name != name {
		w.Name = name
	}
}
func (w *Workflow) UpdateDescription(description string) {
	if w.Description != description {
		w.Description = description
	}
}

func (w *Workflow) AddVersion(param *VersionOption) (*WorkflowVersion, error) {
	applog.Infow("Workflow AddVersion", "param", param)
	if err := param.validate(); err != nil {
		return nil, err
	}
	metaData := make(map[string]string)
	if param.Source == WorkflowSourceGit {
		metaData[WorkflowGitURL] = param.Url
		metaData[WorkflowGitTag] = param.Tag
		//if param.Token != "" {
		//	metaData[WorkflowGitToken] = param.Token
		//}
	}
	// create workflow version
	version := &WorkflowVersion{
		ID:               utils.GenWorkflowVersionID(),
		Status:           WorkflowVersionPendingStatus,
		Message:          "",
		Language:         param.Language,
		MainWorkflowPath: param.MainWorkflowPath,
		Source:           param.Source,
		Inputs:           make([]WorkflowParam, 0),
		Outputs:          make([]WorkflowParam, 0),
		Graph:            "",
		Metadata:         metaData,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// append files
	if w.Versions == nil {
		w.Versions = make(map[string]*WorkflowVersion)
	}
	w.Versions[version.ID] = version
	w.LatestVersion = version.ID
	return version, nil
}

// WorkflowVersion workflow version
type WorkflowVersion struct {
	// ID is the unique identifier of the workflow version
	ID string
	// Status is the status of the workflow version
	Status string
	// Message is the message of the workflow version
	Message string
	// Language is the language of the workflow version
	Language string
	// LanguageVersion is the language version of the workflow version
	LanguageVersion string
	// MainWorkflowPath is the main file of the workflow version
	MainWorkflowPath string
	// Inputs is the inputs of the workflow version
	Inputs []WorkflowParam
	// Outputs is the outputs of the workflow version
	Outputs []WorkflowParam
	// Graph is the graph of the workflow version
	Graph string
	// Metadata is the metadata of the workflow version
	Metadata map[string]string
	// Source is the source of the workflow version, eg. git,file
	Source string
	// Files is the files of the workflow
	Files map[string]*WorkflowFile
	// CreatedAt is the create time of workflow version
	CreatedAt time.Time
	// UpdatedAt is the update time of workflow version
	UpdatedAt time.Time
}

type WorkflowParam struct {
	// Name param name
	Name string `json:"name"`
	// Type param type
	Type string `json:"type"`
	// Optional param is optional
	Optional bool `json:"optional"`
	// Default param default value
	Default string `json:"default,omitempty"`
}

func (v *WorkflowVersion) AddFile(param *FileParam) (*WorkflowFile, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}
	// create file
	file := &WorkflowFile{
		ID:        utils.GenWorkflowFileID(),
		Path:      param.Path,
		Content:   param.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// append files
	if v.Files == nil {
		v.Files = make(map[string]*WorkflowFile)
	}
	v.Files[file.ID] = file

	return file, nil
}

// WorkflowFile workflow file
type WorkflowFile struct {
	// ID is the unique identifier of the workflow file
	ID string
	// Path filename with path
	Path string
	// Content file content
	Content string
	// CreatedAt is the create time of workflow file
	CreatedAt time.Time
	// UpdatedAt is the update time of workflow file
	UpdatedAt time.Time
}
