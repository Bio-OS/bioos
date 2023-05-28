// This file is define ipython notebook format
// see more https://nbformat.readthedocs.io/en/latest/format_description.html
package notebook

import (
	"encoding/json"
	"fmt"

	"github.com/Bio-OS/bioos/pkg/validator"
)

type IPythonNotebook struct {
	Meta        IPythonNotebookMeta   `json:"metadata"`
	Cells       []IPythonNotebookCell `json:"cells"`
	Format      int                   `json:"nbformat" validate:"required,min=4"`
	FormatMinor int                   `json:"nbformat_minor"`
}

type IPythonNotebookMeta struct {
	LanguageInfo IPythonNotebookLanguageInfo `json:"language_info"`
	KernelSpec   IPythonNotebookKernelSpec   `json:"kernelspec"`
	MaxCellID    int                         `json:"max_cell_id"`
}

type IPythonNotebookLanguageInfo struct {
	Name          string `json:"name" validate:"required,oneof=python R"`
	Version       string `json:"version" validate:"required"`
	PygmentsLexer string `json:"pygments_lexer"`
	// codemirror_mode in docs is string but jupyter created is {"name":"ipython",...}
	// CodeMirrorMode string `json:"codemirror_mode"`
	FileExtension string `json:"file_extension" validate:"required"`
	MimeType      string `json:"mimetype" validate:"required"`
}

type IPythonNotebookKernelSpec struct {
	Name        string `json:"name" validate:"required"`
	DisplayName string `json:"display_name"`
	Language    string `json:"language" validate:"required,oneof=python R"`
}

type IPythonNotebookCell struct {
	ID             string                      `json:"id"`
	Type           string                      `json:"cell_type" validate:"required,oneof=code markdown raw"`
	Meta           map[string]interface{}      `json:"metadata"`
	ExecutionCount *int                        `json:"execution_count"`
	Source         []string                    `json:"source"`
	Outputs        []IPythonNotebookCellOutput `json:"outputs"`
	Attachments    map[string]interface{}      `json:"attachments"`
}

type IPythonNotebookCellOutput struct {
	Type string `json:"output_type" validate:"required,oneof=stream display_data execute_result error"`
	Name string `json:"name" validate:"required"`

	Text []string `json:"text"`

	// in display data
	Meta map[string]interface{} `json:"metadata"`
	Data map[string]interface{} `json:"data"`

	// in execute result; used to be pyout / prompt_number
	ExecutionCount int `json:"execute_count"`

	// In errors
	EName     string   `json:"ename"`
	EValue    string   `json:"evalue"`
	Traceback []string `json:"traceback"`
}

func validateIPythonNotebook(blob []byte) error {
	if len(blob) == 0 {
		return fmt.Errorf("empty content")
	}
	var nb IPythonNotebook
	if err := json.Unmarshal(blob, &nb); err != nil {
		return fmt.Errorf("not a valid json string: %w", err)
	}
	return validator.Validate(&nb)
}
