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
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils/bom"
)

// GetDataModelType ...
func GetDataModelType(dataModelName string) string {
	if strings.HasSuffix(dataModelName, consts.DataModelEntitySetNameSuffix) {
		return consts.DataModelTypeEntitySet
	}
	if consts.WorkspaceTypeDataModelName == dataModelName {
		return consts.DataModelTypeWorkspace
	}
	return consts.DataModelTypeEntity
}

// GenDataModelHeaderOfID ...
func GenDataModelHeaderOfID(name string) string {
	return name + "_" + consts.DataModelPrimaryHeader
}

// GenDataModelEntitySetName ...
func GenDataModelEntitySetName(name string) string {
	return name + consts.DataModelEntitySetNameSuffix
}

func ReadDataModelFromCSV(filePath string) ([]string, [][]string, error) {
	if path.Ext(filePath) != ".csv" {
		return nil, nil, fmt.Errorf("%s not support, please use a csv file", path.Ext(filePath))
	}
	var headers []string
	rows := make([][]string, 0)
	// Open the file
	csvRawFile, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't open the csv file: %w", err)
	}
	defer csvRawFile.Close()

	csvfile, _ := bom.ReadSkipBOM(csvRawFile)

	// Parse the file
	r := csv.NewReader(csvfile)

	isHeader := true
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("read csv file failed: %w", err)
		}

		if isHeader {
			headers = record
			isHeader = false
			continue
		}
		rows = append(rows, record)
	}
	return headers, rows, err
}

// WriteDataModelToCSVFile ...
func WriteDataModelToCSVFile(filePath string, headers []string, rows [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileCSV := bom.WriteWithBOM(file, bom.UTF8)

	w := csv.NewWriter(fileCSV)

	// write headers
	err = w.Write(headers)
	if err != nil {
		return err
	}
	w.Flush()

	// write rows
	for _, row := range rows {
		err := w.Write(row)
		if err != nil {
			return err
		}
		// flush buffer
		w.Flush()
	}
	return nil
}
