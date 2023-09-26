//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

// WorkspaceTypedSchema ...
type WorkspaceTypedSchema struct {
	Name        string                 `yaml:"name"`
	Version     string                 `yaml:"version"`
	Description string                 `yaml:"description"`
	DataModels  []DataModelTypedSchema `yaml:"dataModels,omitempty"`
	Workflows   []WorkflowTypedSchema  `yaml:"workflows,omitempty"`
	Notebooks   NotebookTypedSchema    `yaml:"notebooks,omitempty"`
}

// DataModelTypedSchema ...
type DataModelTypedSchema struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

// WorkflowTypedSchema ...
type WorkflowTypedSchema struct {
	Name             string           `yaml:"name"`
	Description      *string          `yaml:"description,omitempty"`
	Language         string           `yaml:"language"`
	Version          *string          `yaml:"version,omitempty"`
	MainWorkflowPath string           `yaml:"mainWorkflowPath"`
	Path             string           `yaml:"path"`
	Metadata         WorkflowMetadata `yaml:"metadata"`
}

// WorkflowMetadata ...
type WorkflowMetadata struct {
	Scheme string  `yaml:"scheme,omitempty"`
	Repo   string  `yaml:"repo,omitempty"`
	Tag    string  `yaml:"tag,omitempty"`
	Token  *string `yaml:"token,omitempty"`
}

// NotebookTypedSchema ...
type NotebookTypedSchema struct {
	Image     *NoteBookImage `yaml:"image,omitempty"`
	Artifacts []*Artifact    `yaml:"artifacts,omitempty"`
}

// NoteBookImage ...
type NoteBookImage struct {
	Name        string   `yaml:"name"`
	DisPlayName string   `yaml:"disPlayName"`
	Packages    string   `yaml:"packages"`
	BasicEnv    []string `yaml:"basicEnv"`
}

// Artifact ...
type Artifact struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
