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

package consts

import "os"

const (
	XRequestIDKey = "X-Request-ID"
)

const (
	QuerySliceDelimiter = ","
	OrderDelimiter      = ":"
	ASCOrdering         = "asc"
	DESCOrdering        = "desc"
)

// DataModel's restriction
const (
	DataModelEntitySetNameSuffix = "_set"
	WorkspaceTypeDataModelName   = "workspace_data"

	DataModelTypeEntity    = "entity"
	DataModelTypeEntitySet = "entity_set"
	DataModelTypeWorkspace = "workspace"

	WorkspaceTypeDataModelHeaderKey    = "Key"
	WorkspaceTypeDataModelHeaderValue  = "Value"
	WorkspaceTypeDataModelMaxHeaderNum = 2
	DataModelPrimaryHeader             = "id"

	DataModelRefPrefix              = "this."
	WorkspaceTypeDataModelRefPrefix = "workspace."
)

// Submission's type
const (
	DataModelTypeSubmission = "dataModel"
	FilePathTypeSubmission  = "filePath"
)

const (
	SubmissionPending    = "Pending"
	SubmissionRunning    = "Running"
	SubmissionFailed     = "Failed"
	SubmissionFinished   = "Finished"
	SubmissionSucceeded  = "Succeeded"
	SubmissionCancelling = "Cancelling"
	SubmissionCancelled  = "Cancelled"
)

// status enum for run
const (
	RunPending    = "Pending"
	RunRunning    = "Running"
	RunFailed     = "Failed"
	RunSucceeded  = "Succeeded"
	RunCancelling = "Cancelling"
	RunCancelled  = "Cancelled"
)

// status enum for task
const (
	TaskQueued       = "Queued"
	TaskInitializing = "Initializing"
	TaskRunning      = "Running"
	TaskFailed       = "Failed"
	TaskSucceeded    = "Succeeded"
	TaskCancelled    = "Cancelled"
)

// submission status groups
var (
	NonFinishedSubmissionStatuses = []string{SubmissionPending, SubmissionRunning, SubmissionCancelling}
	FinishedSubmissionStatuses    = []string{SubmissionCancelled, SubmissionFinished}
	AllowCancelSubmissionStatuses = []string{SubmissionPending, SubmissionRunning}
)

// run status groups
var (
	NonFinishedRunStatuses = []string{RunPending, RunRunning, RunCancelling}
	FinishedRunStatuses    = []string{RunCancelled, RunSucceeded, RunFailed}
	AllowCancelRunStatuses = []string{RunPending, RunRunning}
)

// task status groups
var (
	NonFinishedTaskStatuses = []string{TaskQueued, TaskInitializing, TaskRunning}
	FinishedTaskStatuses    = []string{TaskSucceeded, TaskFailed, TaskCancelled}
)

// workspace datamodel value header
const WsDataModelValueHeader = "Value"

const ImportWorkspaceFileTypeExt = ".zip"

const (
	WorkspaceDir                 = "workspace"
	WorkspaceCoverImageName      = "cover.png"
	WorkspaceYAMLName            = "workspace.yaml"
	WorkspaceZIPName             = "workspace.zip"
	WorkflowDirName              = "workflow"
	NotebookDirName              = "notebook"
	FileDirName                  = "file"
	SubmissionDirName            = "submission"
	DataModelDirName             = "data"
	WorkspaceScopedSchemaVersion = "1"
)

// FileMode ...
const (
	SchemaFileMode = os.FileMode(0777)
)

// WorkflowType ...
const (
	WorkflowTypeWDL      = "WDL"
	WorkflowTypeNextflow = "Nextflow"
)
