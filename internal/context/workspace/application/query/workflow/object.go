package workflow

import "time"

type Workflow struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	WorkspaceID   string           `json:"workflowID"`
	Description   string           `json:"description"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	LatestVersion *WorkflowVersion `json:"latestVersion"`
}

type WorkflowVersion struct {
	ID               string              `json:"id"`
	Status           string              `json:"status"`
	Message          string              `json:"message"`
	Language         string              `json:"language"`
	LanguageVersion  string              `json:"languageVersion"`
	MainWorkflowPath string              `json:"mainWorkflowPath"`
	Inputs           []WorkflowParam     `json:"inputs,omitempty"`
	Outputs          []WorkflowParam     `json:"outputs,omitempty"`
	Graph            string              `json:"graph,omitempty"`
	Source           string              `json:"source"`
	Files            []*WorkflowFileInfo `json:"files"`
	Metadata         map[string]string   `json:"metadata,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

type WorkflowParam struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Default  string `json:"default,omitempty"`
}

type WorkflowFile struct {
	ID                string    `json:"id"`
	WorkflowVersionID string    `json:"workflowVersionID"`
	Path              string    `json:"path"`
	Content           string    `json:"content"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (f *WorkflowFile) ToWorkflowFileInfo() *WorkflowFileInfo {
	return &WorkflowFileInfo{
		ID:   f.ID,
		Path: f.Path,
	}
}

type WorkflowFileInfo struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}
