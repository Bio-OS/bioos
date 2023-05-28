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
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	MAX_COUNT       = 10000
	MAX_SIZE        = 1 * 1024 * 1024 * 1024
	MAX_PATH_LENGTH = 512
)

// ValidateFSDirectory check specified path exist and is folder
// return os.ErrNotExist if not exist
func ValidateFSDirectory(dirname string) error {
	if !fs.ValidPath(strings.Trim(dirname, "/")) {
		return fmt.Errorf("path '%s' invalid", dirname)
	}
	f, err := os.Open(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("can not open dir '%s': %w", dirname, err)
	}
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("get file '%s' stat fail: %w", dirname, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("file '%s' is not a directory", dirname)
	}
	return nil
}

// ValidateFileExist check if given file exist in path.
func ValidateFileExist(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file: %s does not exist", filePath)
		}
		return err
	}
	return nil
}

// GetSubPath check if target path is subpath of basepath and return subpath
func GetSubPath(basepath, targpath string) (string, bool) {
	rel, err := filepath.Rel(basepath, targpath)
	// if no any relative, the error will be:
	//   errors.New("Rel: can't make " + targpath + " relative to " + basepath)
	if err == nil && !strings.HasPrefix(rel, "..") {
		return rel, true
	}
	return "", false
}

func Unzip(zipFilePath, targetDir string) error {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	var count uint64 = 0
	var size uint64 = 0
	// Iterate through each file in the zip archive
	for _, zipFile := range zipReader.File {
		err = func() error {
			// Open the compressed file
			rc, err := zipFile.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			count += 1
			if count > MAX_COUNT {
				return fmt.Errorf("too many files")
			}
			size += zipFile.UncompressedSize64
			if size > MAX_SIZE {
				return fmt.Errorf("files total size over 1GB")
			}

			// Create the target file path
			targetFilePath := filepath.Join(targetDir, zipFile.Name)
			// If target file path is not valid, return an error
			if !strings.HasPrefix(targetFilePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
				return fmt.Errorf("%s: illegal file path", targetFilePath)
			}
			// If target file path is too long, return an error
			if len(targetFilePath) > MAX_PATH_LENGTH {
				return fmt.Errorf("%s: file path too long", targetFilePath)
			}

			// If the file is a directory, create it
			if zipFile.FileInfo().IsDir() {
				os.MkdirAll(targetFilePath, zipFile.Mode())
			} else {
				// If the file is not a directory, decompress it
				if err = os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm); err != nil {
					return err
				}
				// Open the target file for writing
				outFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
				if err != nil {
					return err
				}
				defer outFile.Close()
				//TODO Decompress the file using zlib
				//zr, err := zlib.NewReader(rc)
				//if err != nil {
				//	return err
				//}
				//defer zr.Close()

				// Write the decompressed data to the target file
				if _, err = io.Copy(outFile, rc); err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipDir(srcDir string, zipFileName string) error {
	// remove the old files
	os.RemoveAll(zipFileName)

	zipFile, err := os.Create(filepath.Clean(zipFileName))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// zip all path in dir
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {

		if path == srcDir {
			return nil
		}

		// get zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, srcDir+`/`)

		if info.IsDir() {
			header.Name += `/`
		} else {
			// set zip algorithm
			header.Method = zip.Deflate
		}

		// create zip header
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(filepath.Clean(path))
			if err != nil {
				return err
			}
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
	return err
}
