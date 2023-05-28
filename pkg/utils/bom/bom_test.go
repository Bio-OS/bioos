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

package bom

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/onsi/gomega"
)

var (
	rowHeaders = []string{"111", "222"}
	rowRows    = [][]string{{"111", "222"}}
	encodings  = []Encoding{UTF8, UTF16BigEndian, UTF16LittleEndian, UTF32BigEndian, UTF32LittleEndian}
)

func TestBOM(t *testing.T) {
	g := gomega.NewWithT(t)
	fileDir := os.TempDir()
	err := os.MkdirAll(fileDir, os.FileMode(0777))
	g.Expect(err).To(gomega.BeNil())
	filePath := path.Join(os.TempDir(), "test.csv")
	for _, encoding := range encodings {
		err = createFile(filePath, encoding)
		g.Expect(err).To(gomega.BeNil())
		err = readFile(filePath, encoding)
		g.Expect(err).To(gomega.BeNil())
		err = os.RemoveAll(filePath)
		g.Expect(err).To(gomega.BeNil())
	}
}

func createFile(filePath string, e Encoding) error {
	fs, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fs.Close()
	writer := WriteWithBOM(fs, e)
	w := csv.NewWriter(writer)

	// write headers
	err = w.Write(rowHeaders)
	if err != nil {
		return err
	}
	w.Flush()

	// write rows
	for _, row := range rowRows {
		err := w.Write(row)
		if err != nil {
			return err
		} // flush buffer
		w.Flush()
	}
	return nil
}

func readFile(filePath string, e Encoding) error {
	fs, err := os.Open(filePath)
	if err != nil {
		return err
	}
	fileCSV, encoding := ReadSkipBOM(fs)
	defer fs.Close()
	if encoding != e {
		return fmt.Errorf("detected encoding[%s] is not %s ", encoding.String(), e.String())
	}
	r := csv.NewReader(fileCSV)
	headers, err := r.Read()
	rows := make([][]string, 0)
	if err != nil {
		return err
	}
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		rows = append(rows, row)
	}
	if !reflect.DeepEqual(headers, rowHeaders) || !reflect.DeepEqual(rows, rowRows) {
		return fmt.Errorf("headers or rows not equal! ")
	}
	return nil
}
