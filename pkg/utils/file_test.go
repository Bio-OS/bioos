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
	"os"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestValidateFSDirectory(t *testing.T) {
	// prepare file
	prefix := os.TempDir()
	regularFile, err := os.CreateTemp(prefix, "TestValidateFSDirectory_temp_")
	if err != nil {
		t.Fatalf("can not create temp file: %s", err)
	}
	defer os.Remove(regularFile.Name())
	t.Logf("created temp file: %s", regularFile.Name())

	cases := []struct {
		dirname     string
		expectError bool
	}{
		{prefix, false},
		{regularFile.Name(), true},
		{prefix + "/path/not/exist", true},
		{"/usr/sbin/", false},
	}
	g := gomega.NewWithT(t)
	for _, c := range cases {
		var match types.GomegaMatcher
		if c.expectError {
			match = gomega.HaveOccurred()
		} else {
			match = gomega.BeNil()
		}
		g.Expect(ValidateFSDirectory(c.dirname)).To(match)
	}
}

func TestIsSubPath(t *testing.T) {
	g := gomega.NewWithT(t)
	base := "/base"
	cases := []struct {
		target  string
		subpath string
		expect  bool
	}{
		{"/base/haha", "haha", true},
		{"/other/haha", "", false},
		{"./base/haha", "", false},
	}
	for _, c := range cases {
		s, ok := GetSubPath(base, c.target)
		g.Expect(ok).To(gomega.Equal(c.expect))
		if ok {
			g.Expect(s).To(gomega.Equal(c.subpath))
		}
	}
}
