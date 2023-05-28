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

package utils

import (
	"fmt"
	"strings"

	"github.com/rs/xid"
	"k8s.io/apimachinery/pkg/util/rand"
)

// CharRange help to specify the range of character for function randString.
type charRange string

const (
	// UpperLetterRange all upper letters.
	upperLetterRange charRange = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// NumberRange all numbers.
	numberRange charRange = "0123456789"
)

func genResourceID(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, xid.New().String())
}

func randString(n int, charRange charRange) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charRange[rand.Intn(len(charRange)-1)])
	}
	return sb.String()
}

// GenPublicWorkspaceID ...
func GenPublicWorkspaceID() string {
	return randString(3, upperLetterRange) + randString(6, numberRange)
}

// GenWorkspaceID ...
func GenWorkspaceID() string {
	return genResourceID("w")
}

// GenWorkflowID ...
func GenWorkflowID() string {
	return genResourceID("f")
}

// GenWorkflowVersionID ...
func GenWorkflowVersionID() string {
	return genResourceID("v")
}

// GenWorkflowFileID ...
func GenWorkflowFileID() string {
	return genResourceID("wf")
}

// GenSubmissionID ...
func GenSubmissionID() string {
	return genResourceID("s")
}

// GenRunID ...
func GenRunID() string {
	return genResourceID("r")
}

// GenNotebookServerID ...
func GenNotebookServerID() string {
	return genResourceID("n")
}

// GenClusterID ...
func GenClusterID() string {
	return genResourceID("c")
}

// GenDataModelID ...
func GenDataModelID() string {
	return genResourceID("d")
}
