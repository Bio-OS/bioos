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

package validator

import "regexp"

// length validator.
const (
	MinResNameLength = 1
	MaxResNameLength = 200

	MinWorkspaceDescLength = 0
	MaxWorkspaceDescLength = 1000

	MinDataModelNameLength          = 1
	MaxDataModelNameLength          = 30
	MaxEntitySetDataModelNameLength = 50
	MaxDataModelHeaderLength        = 100
	MaxDataModelRowIDLength         = 100

	MinSubmissionNameSuffixLength = 1
	MaxSubmissionNameSuffixLength = 200
	// SubmissionName = WorkflowName + "-history-" + Suffix
	MinSubmissionNameLength = MinResNameLength + 9 + MinSubmissionNameSuffixLength
	MaxSubmissionNameLength = MaxResNameLength + 9 + MaxSubmissionNameSuffixLength
	MaxSubmissionDescLength = 1000
)

var (
	ResNameRegex = regexp.MustCompile(`^[\p{Han}A-Za-z0-9][-_\p{Han}A-Za-z0-9]*$`)
	// DataModelNameReg _${data_model}(data model name) is the data model stored for submission
	DataModelNameReg = regexp.MustCompile("^[0-9a-zA-Z_][0-9a-zA-Z-_]*$")
	// DataModelHeaderReg _${data_model}_id(data model id header) is the data model stored for submission
	DataModelHeaderReg = regexp.MustCompile("^[0-9a-zA-Z_][0-9a-zA-Z-_]*$")
)
