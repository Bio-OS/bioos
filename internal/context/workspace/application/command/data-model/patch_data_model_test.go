package datamodel

import (
	"reflect"
	"testing"

	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/onsi/gomega"
)

func TestGetDataModelType(t *testing.T) {
	g := gomega.NewWithT(t)

	testCases := []struct {
		describe      string
		dataModelName string
		dataModelType string
	}{
		{
			describe:      "workspace",
			dataModelName: "workspace_data",
			dataModelType: consts.DataModelTypeWorkspace,
		},
		{
			describe:      "entity",
			dataModelName: "12Abc-test",
			dataModelType: consts.DataModelTypeEntity,
		},
		{
			describe:      "entity_set",
			dataModelName: "12Abc-test_set",
			dataModelType: consts.DataModelTypeEntitySet,
		},
		{
			describe:      "entity_set_set",
			dataModelName: "12Abc-test_set_set",
			dataModelType: consts.DataModelTypeEntitySet,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			g.Expect(utils.GetDataModelType(tc.dataModelName)).To(gomega.Equal(tc.dataModelType))
		})
	}
}

func TestRows2Columns(t *testing.T) {
	g := gomega.NewWithT(t)

	headers := []string{"1", "2", "3", "4"}
	rows := [][]string{
		{"5", "6", "7", "8"},
		{"9", "10", "11", "12"},
		{"13", "14", "15", "16"},
	}
	columns := map[string][]string{
		"1": {"5", "9", "13"},
		"2": {"6", "10", "14"},
		"3": {"7", "11", "15"},
		"4": {"8", "12", "16"},
	}

	res := rows2Columns(headers, rows)
	for header, column := range res {
		g.Expect(reflect.DeepEqual(columns[header], column)).To(gomega.BeTrue())
	}
}

func TestColumns2Rows(t *testing.T) {
	g := gomega.NewWithT(t)

	headers := []string{"1", "2", "3", "4"}
	rows := [][]string{
		{"5", "6", "7", "8"},
		{"9", "10", "11", "12"},
		{"13", "14", "15", "16"},
	}
	columns := map[string][]string{
		"1": {"5", "9", "13"},
		"2": {"6", "10", "14"},
		"3": {"7", "11", "15"},
		"4": {"8", "12", "16"},
	}

	res := columns2Rows(headers, columns)
	for index, row := range res {
		g.Expect(reflect.DeepEqual(row, rows[index])).To(gomega.BeTrue())
	}
}

func TestGetRowIDs(t *testing.T) {
	g := gomega.NewWithT(t)

	rows := [][]string{
		{"1", "2", "3", "4"},
		{"5", "6", "7", "8"},
		{"9", "10", "11", "12"},
		{"13", "14", "15", "16"},
	}
	rowIDs := []string{"1", "5", "9", "13"}

	res := getRowIDs(rows)
	g.Expect(reflect.DeepEqual(res, rowIDs)).To(gomega.BeTrue())
}

func TestGenNewHeadersAndRows(t *testing.T) {
	g := gomega.NewWithT(t)

	testCases := []struct {
		describe    string
		dbHeaders   []string
		fileHeaders []string
		dbColumns   map[string][]string
		fileRows    [][]string
		headers     []string
		rows        [][]string
	}{
		{
			describe:    "normal",
			dbHeaders:   []string{"1", "2", "3"},
			fileHeaders: []string{"1", "3", "4", "5"},
			dbColumns: map[string][]string{
				"1": {"6", "11", ""},
				"2": {"7", "12", ""},
				"3": {"8", "", ""},
			},
			fileRows: [][]string{
				{"6", "8", "9", "10"},
				{"11", "13", "14", "15"},
				{"16", "18", "19", "20"},
			},
			headers: []string{"1", "2", "3", "4", "5"},
			rows: [][]string{
				{"6", "7", "8", "9", "10"},
				{"11", "12", "13", "14", "15"},
				{"16", "", "18", "19", "20"},
			},
		},
		{
			describe:    "normal2",
			dbHeaders:   []string{"1", "2"},
			fileHeaders: []string{"1", "3"},
			dbColumns: map[string][]string{
				"1": {"111", "333"},
				"2": {"222", "444"},
			},
			fileRows: [][]string{
				{"111", "555"},
				{"333", "666"},
			},
			headers: []string{"1", "2", "3"},
			rows: [][]string{
				{"111", "222", "555"},
				{"333", "444", "666"},
			},
		},
		{
			describe:    "normal2",
			dbHeaders:   []string{"1", "2", "3"},
			fileHeaders: []string{"1", "3"},
			dbColumns: map[string][]string{
				"1": {""},
				"2": {""},
				"3": {""},
			},
			fileRows: [][]string{
				{"222", "555"},
			},
			headers: []string{"1", "2", "3"},
			rows: [][]string{
				{"222", "", "555"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.describe, func(t *testing.T) {
			headers, rows := genNewHeadersAndRows(tc.dbHeaders, tc.fileHeaders, tc.dbColumns, tc.fileRows)
			g.Expect(reflect.DeepEqual(headers, tc.headers)).To(gomega.BeTrue())
			for index, row := range rows {
				g.Expect(reflect.DeepEqual(row, tc.rows[index])).To(gomega.BeTrue())
			}
		})
	}
}
