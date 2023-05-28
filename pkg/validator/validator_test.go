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

import (
	"math/rand"
	"os"
	"testing"

	"github.com/Bio-OS/bioos/pkg/consts"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"

	"github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	applog.RegisterLogger(&applog.Options{
		Level: "fatal",
	})
	os.Exit(m.Run())
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_你好我是")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	b[0] = 'a'
	return string(b)
}

func TestValidateResName(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Name string `validate:"resName"`
	}

	testCases := []struct {
		describe string
		name     string
		expMatch bool
	}{
		{
			describe: "normal",
			name:     "12Abc-test-测试",
			expMatch: true,
		},
		{
			describe: "normal _",
			name:     "12Abc_test_测试",
			expMatch: true,
		},
		{
			describe: "normal mix",
			name:     "12Abc-test__-测试_test-test12",
			expMatch: true,
		},
		{
			describe: "too short",
			name:     "",
			expMatch: false,
		},
		{
			describe: "normal length",
			name:     randStringRunes(31),
			expMatch: true,
		},
		{
			describe: "too long",
			name:     randStringRunes(201),
			expMatch: false,
		},
		{
			describe: "invalid with other symbol",
			name:     "12Abc_test%测试",
			expMatch: false,
		},
	}

	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{Name: tc.name})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateWorkspaceDesc(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Desc string `validate:"workspaceDesc"`
	}

	testCases := []struct {
		describe string
		desc     string
		expMatch bool
	}{
		{
			describe: "normal",
			desc:     "测试",
			expMatch: true,
		},
		{
			describe: "too short",
			desc:     "",
			expMatch: false,
		},
		{
			describe: "normal long length",
			desc: "测试testABCD测试testABCD测试testABCD测试testABCD测试testABCD测试testABCD测试testABCD测试testABCD" +
				"测试testABCD测试testABCDx",
			expMatch: true,
		},
		{
			describe: "too long",
			desc:     randStringRunes(1001),
			expMatch: false,
		},
	}

	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{Desc: tc.desc})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateNFSMountPath(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Path string `validate:"nfsMountPath"`
	}

	testCases := []struct {
		describe string
		p        string
		expMatch bool
	}{
		{
			describe: "normal",
			p:        "/abc/ddd",
			expMatch: true,
		},
		{
			describe: "wrong",
			p:        "abc/ddd",
			expMatch: false,
		},
	}

	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{Path: tc.p})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateDataModelName(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Name string `validate:"dataModelName"`
	}

	testCases := []struct {
		describe string
		name     string
		expMatch bool
	}{
		{
			describe: "normal",
			name:     "12Abc-test",
			expMatch: true,
		},
		{
			describe: "too short",
			name:     "",
			expMatch: false,
		},
		{
			describe: "too long",
			name:     "testtesttesttesttesttesttesttes",
			expMatch: false,
		},
		{
			describe: "valid with '_'",
			name:     "abc_d",
			expMatch: true,
		},
		{
			describe: "valid with '_' at prefix",
			name:     "_12Abc-t_est",
			expMatch: true,
		},
		{
			describe: "invalid",
			name:     "abc&d",
			expMatch: false,
		},
		{
			describe: "set name up to 50",
			name:     "t0123456789012345678901234567890123456789_set",
			expMatch: true,
		},
		{
			describe: "set name exceed 50",
			name:     "t01234567890123456789012345678901234567890123456789_set",
			expMatch: false,
		},
	}

	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{Name: tc.name})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateDataModelHeaders(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Headers []string `validate:"dataModelHeaders"`
		Name    string
	}

	testCases := []struct {
		describe      string
		headers       []string
		dataModelName string
		expMatch      bool
	}{
		{
			describe:      "normal",
			headers:       []string{"test_id", "xxx-yyy", "xxx-yyy"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "workspace_type",
			headers:       []string{consts.WorkspaceTypeDataModelHeaderKey, consts.WorkspaceTypeDataModelHeaderValue},
			dataModelName: consts.WorkspaceTypeDataModelName,
			expMatch:      true,
		},
		{
			describe:      "workspace_type_with_wrong_headers_number",
			headers:       []string{consts.WorkspaceTypeDataModelHeaderKey, consts.WorkspaceTypeDataModelHeaderValue, "value"},
			dataModelName: consts.WorkspaceTypeDataModelName,
			expMatch:      false,
		},
		{
			describe:      "workspace_type_with_wrong_headers_value",
			headers:       []string{consts.WorkspaceTypeDataModelHeaderKey, "value"},
			dataModelName: consts.WorkspaceTypeDataModelName,
			expMatch:      false,
		},
		{
			describe:      "workspace_type_with_wrong_headers_value",
			headers:       []string{utils.GenDataModelHeaderOfID(consts.WorkspaceTypeDataModelName), "value"},
			dataModelName: consts.WorkspaceTypeDataModelName,
			expMatch:      false,
		},
		{
			describe:      "only id",
			headers:       []string{"test_id"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "match with single word",
			headers:       []string{"test_id", "a"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "match with 20 headers",
			headers:       []string{"test_id", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "match with 'xxx_yyy' headers",
			headers:       []string{"test_id", "0_1_3"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "match with '_xxx_yyy' headers",
			headers:       []string{"_test_id", "0_1_3"},
			dataModelName: "_test",
			expMatch:      true,
		},
		{
			describe:      "not match with 'doc_id' headers when dataModelName is 'test'",
			headers:       []string{"doc_id", "_1"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with 'id' headers",
			headers:       []string{"id", "_1"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with '_yyy' headers",
			headers:       []string{"test_id", "_1"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "match set headers",
			headers:       []string{"test_set_id", "test"},
			dataModelName: "test_set",
			expMatch:      true,
		},
		{
			describe:      "not match set headers",
			headers:       []string{"test_set_id", "test_"},
			dataModelName: "test_set",
			expMatch:      false,
		},
		{
			describe:      "not match with 100 characters",
			headers:       []string{"test_id", "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"},
			dataModelName: "test",
			expMatch:      true,
		},
		{
			describe:      "not match with 51 headers",
			headers:       []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with -xxx",
			headers:       []string{"test_id", "-aaaaaa"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with -",
			headers:       []string{"test_id", "-"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with empty list",
			headers:       []string{},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with chinese words",
			headers:       []string{"test_id", "测试"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with specific symbol",
			headers:       []string{"test_id", "!"},
			dataModelName: "test",
			expMatch:      false,
		},
		{
			describe:      "not match with more than 100 characters",
			headers:       []string{"test_id", "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"},
			dataModelName: "test",
			expMatch:      false,
		},
	}
	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err = Validate(Obj{Headers: tc.headers, Name: tc.dataModelName})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateDeleteDataModelHeaders(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Headers []string `validate:"deleteDataModelHeaders"`
	}

	testCases := []struct {
		describe string
		headers  []string
		expMatch bool
	}{
		{
			describe: "normal",
			headers:  []string{"xxx-yyy"},
			expMatch: true,
		},
		{
			describe: "match with empty list",
			headers:  []string{},
			expMatch: true,
		},
		{
			describe: "not match with empty string",
			headers:  []string{""},
			expMatch: false,
		},
		{
			describe: "not match with not only empty string",
			headers:  []string{"xxx-yyy", ""},
			expMatch: false,
		},
	}
	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{tc.headers})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateDataModelRows(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Name string
		Rows [][]string `validate:"dataModelRows"`
	}

	testCases := []struct {
		describe string
		name     string
		rows     [][]string
		expMatch bool
	}{
		{
			describe: "normal",
			rows:     [][]string{{"xxx-yyy"}},
			expMatch: true,
		},
		{
			describe: "normal two row",
			rows:     [][]string{{"xxx-yyy"}, {"aaa-bbb"}},
			expMatch: true,
		},
		{
			describe: "match with -xxx",
			rows:     [][]string{{"-"}},
			expMatch: true,
		},
		{
			describe: "not match with chinese words",
			rows:     [][]string{{"测试"}},
			expMatch: true,
		},
		{
			describe: "match with specific symbol",
			rows:     [][]string{{"!"}},
			expMatch: true,
		},
		{
			describe: "match with not only \"\"",
			rows:     [][]string{{"column2", ""}},
			expMatch: true,
		},
		{
			describe: "not match with empty list",
			rows:     [][]string{},
			expMatch: false,
		},
		{
			describe: "not match with empty list",
			rows:     [][]string{{}},
			expMatch: false,
		},
		{
			describe: "not match with one \"\" row",
			rows:     [][]string{{""}},
			expMatch: false,
		},
		{
			describe: "not match with two \"\" rows",
			rows:     [][]string{{""}, {""}},
			expMatch: false,
		},
		{
			describe: "not match with \"\" row and not empty row",
			rows: [][]string{
				{"1111", "222222", "333333"},
				{""},
			},
			expMatch: false,
		},
		{
			describe: "not match with \"\" & not empty grid and one not empty row",
			rows: [][]string{
				{"1111", "222222", "333333"},
				{"", "11111"},
			},
			expMatch: false,
		},
		{
			describe: "not match with only \"\" and not empty grid",
			rows: [][]string{
				{"", "11111"},
			},
			expMatch: false,
		},
		{
			describe: "match set row",
			name:     "h_set",
			rows:     [][]string{{"xxx-yyy", `["a","b"]`}, {"aaa-bbb", `["c"]`}},
			expMatch: true,
		},
		{
			describe: "not match set row count",
			name:     "h_set",
			rows:     [][]string{{"xxx-yyy"}, {"aaa-bbb", `["c"]`}},
			expMatch: false,
		},
		{
			describe: "not match set row content",
			name:     "h_set",
			rows:     [][]string{{"xxx-yyy", `["a","b"`}, {"aaa-bbb", `["c"]`}},
			expMatch: false,
		},
	}
	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err = Validate(Obj{tc.name, tc.rows})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateSubmissionName(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Name string `validate:"submissionName"`
	}

	testCases := []struct {
		describe string
		name     string
		expMatch bool
	}{
		{
			describe: "normal",
			name:     "测试pp-history-abcd12",
			expMatch: true,
		},
		{
			describe: "too short",
			name:     "a-history-",
			expMatch: false,
		},
		{
			describe: "too long",
			name:     randStringRunes(410),
			expMatch: false,
		},
		{
			describe: "invalid",
			name:     "abcd-测试_t",
			expMatch: false,
		},
	}

	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{Name: tc.name})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}

func TestValidateSubmissionDesc(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Description *string `validate:"omitempty,submissionDesc"`
	}

	testCases := []struct {
		describe string
		desc     *string
		expMatch bool
	}{
		{
			describe: "normal",
			desc:     utils.PointString("测试"),
			expMatch: true,
		},
		{
			describe: "normal empty",
			desc:     utils.PointString(""),
			expMatch: true,
		},
		{
			describe: "nil",
			desc:     nil,
			expMatch: true,
		},
		{
			describe: "too long",
			desc:     utils.PointString(randStringRunes(1001)),
			expMatch: false,
		},
	}

	err := RegisterValidators()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			err := Validate(Obj{Description: tc.desc})
			g.Expect(tc.expMatch).To(gomega.Equal(err == nil))
		})
	}
}
