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
	"regexp"
	"testing"

	. "github.com/onsi/gomega"
)

var UpperLetterRegexp = regexp.MustCompile(`[a-zA-Z]+`)
var NumberRegexp = regexp.MustCompile(`[0-9]+`)

func TestGenPublicWorkspaceID(t *testing.T) {
	g := NewWithT(t)

	for i := 0; i < 100; i++ {
		res := GenPublicWorkspaceID()
		if !g.Expect(len(res)).To(BeEquivalentTo(9)) {
			t.Fatalf("the length of the result string should be 9: %s", res)
		}
		isFirstPartUpper := UpperLetterRegexp.MatchString(res[0:3])
		isSecondPartNumber := NumberRegexp.MatchString(res[3:])
		if !g.Expect(isFirstPartUpper).To(BeTrue()) {
			t.Fatalf("the first part of the result string should be three upper letters: %s", res)
		}
		if !g.Expect(isSecondPartNumber).To(BeTrue()) {
			t.Fatalf("the second part of the result string should be six upper letters: %s", res)
		}
	}
}
