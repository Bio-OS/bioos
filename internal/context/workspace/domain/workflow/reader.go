package workflow

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils/exec"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ReaderOptions workflow reader options
type ReaderOptions struct {
	WomtoolPath string
}

// Reader workflow reader
type Reader interface {
	ParseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error)
	ValidateWorkflowFiles(ctx context.Context, version *WorkflowVersion, baseDir, mainWorkflowPath string) error
	GetWorkflowInputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error)
	GetWorkflowOutputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error)
	GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error)
}

// WDLReader reader for workflow written by WDL
type WDLReader struct {
	options *ReaderOptions
}

// ParseWorkflowVersion ...
func (r *WDLReader) ParseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error) {
	versionRegexp := regexp.MustCompile(VersionRegexpStr)
	file, err := os.Open(mainWorkflowPath)
	if err != nil {
		applog.Errorw("fail to open main workflow file", "err", err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matched := versionRegexp.FindStringSubmatch(line)
		if matched != nil && len(matched) >= 2 {
			applog.Infow("version regexp matched", "matched", matched)
			return matched[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "draft-2", nil
}

// ValidateWorkflowFiles ...
func (r *WDLReader) ValidateWorkflowFiles(ctx context.Context, version *WorkflowVersion, baseDir, mainWorkflowPath string) error {
	applog.Infow("start to validate files", "mainWorkflowPath", mainWorkflowPath)
	validateResult, err := exec.Exec(ctx, CommandExecuteTimeout, "java", "-jar", r.options.WomtoolPath, "validate", path.Join(baseDir, mainWorkflowPath), "-l")
	if err != nil {
		applog.Errorw("fail to validate workflow", "err", err, "result", string(validateResult))
		return apperrors.NewInternalError(fmt.Errorf("validate workflow failed"))
	}
	validateResultLines := strings.Split(string(validateResult), "\n")
	applog.Infow("validate result", "result", validateResultLines)
	// parse and save workflow files
	if len(validateResultLines) < 2 || strings.ToLower(validateResultLines[0]) != "success!" {
		return proto.ErrorWorkflowValidateError("fail to validate workflow version:%s", version.ID)
	}
	workflowFiles := []string{mainWorkflowPath}
	// need to start from line 2(start with line 0)
	for i := 2; i < len(validateResultLines); i++ {
		absPath := validateResultLines[i]
		if len(absPath) == 0 {
			continue
		}
		// validate file
		if _, err := os.Stat(absPath); err == nil {
			// in mac absPath was prefix with /private
			relPath, err := filepath.Rel(baseDir, absPath[strings.LastIndex(absPath, baseDir):])
			if err != nil {
				return apperrors.NewInternalError(err)
			}
			applog.Infow("file path", "baseDir", baseDir, "absPath", absPath, "relPath", relPath)
			workflowFiles = append(workflowFiles, relPath)
		}
	}
	for _, relPath := range workflowFiles {
		input, err := os.ReadFile(path.Join(baseDir, relPath))
		if err != nil {
			applog.Errorw("fail to read file", "err", err)
			return apperrors.NewInternalError(err)
		}

		encodedContent := base64.StdEncoding.EncodeToString(input)

		workflowFile, err := version.AddFile(&FileParam{
			Path:    relPath,
			Content: encodedContent,
		})
		if err != nil {
			return err
		}
		applog.Infow("success add workflow file", "workflowVersionID", version.ID, "fileID", workflowFile.ID, "path", workflowFile.Path)
	}
	return nil
}

// GetWorkflowInputs ...
func (r *WDLReader) GetWorkflowInputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error) {
	return r.getWorkflowParams(ctx, "java", "-jar", r.options.WomtoolPath, "inputs", WorkflowFilePath)
}

// GetWorkflowOutputs ...
func (r *WDLReader) GetWorkflowOutputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error) {
	return r.getWorkflowParams(ctx, "java", "-jar", r.options.WomtoolPath, "outputs", WorkflowFilePath)
}

// GetWorkflowGraph ...
func (r *WDLReader) GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error) {
	graph, err := exec.Exec(ctx, CommandExecuteTimeout, "java", "-jar", r.options.WomtoolPath, "graph", WorkflowFilePath)
	if err != nil {
		return "", err
	}

	return string(graph), nil
}

// getWorkflowParams ...
func (r *WDLReader) getWorkflowParams(ctx context.Context, name string, arg ...string) ([]WorkflowParam, error) {
	params := make([]WorkflowParam, 0)
	outputsResult, err := exec.Exec(ctx, CommandExecuteTimeout, name, arg...)
	if err != nil {
		return params, err
	}
	var outputsMap map[string]string
	if err := json.Unmarshal(outputsResult, &outputsMap); err != nil {
		return params, err
	}

	for paramName, value := range outputsMap {
		paramType, optional, defaultValue := parseWorkflowParamValue(value)
		param := WorkflowParam{
			Name:     paramName,
			Type:     paramType,
			Optional: optional,
		}
		if defaultValue != nil {
			param.Default = *defaultValue
		}
		params = append(params, param)
	}
	// keep the return sort stable
	sort.Slice(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})
	return params, nil
}

func parseWorkflowParamValue(value string) (paramType string, optional bool, defaultValue *string) {
	splitByLeftBracket := strings.SplitN(value, "(", 2)
	paramType = strings.TrimSpace(splitByLeftBracket[0])
	if len(splitByLeftBracket) == 1 {
		return paramType, false, nil
	}

	extraInfo := strings.TrimSuffix(splitByLeftBracket[1], ")")
	splitByComma := strings.SplitN(extraInfo, ",", 2)
	if strings.TrimSpace(splitByComma[0]) == "optional" {
		optional = true
	}
	if len(splitByComma) == 1 {
		return paramType, optional, nil
	}

	defaultInfo := strings.TrimSpace(splitByComma[1])
	splitByEqual := strings.SplitN(defaultInfo, "=", 2)
	if len(splitByEqual) != 2 || strings.ToLower(strings.TrimSpace(splitByEqual[0])) != "default" {
		return paramType, optional, nil
	}
	rawDefaultValue := strings.TrimSpace(splitByEqual[1])
	defaultValue = new(string)
	if strings.HasPrefix(rawDefaultValue, `"`) && strings.HasSuffix(rawDefaultValue, `"`) { // String type
		_ = json.Unmarshal([]byte(rawDefaultValue), defaultValue) // escape, never error
	} else {
		*defaultValue = rawDefaultValue
	}
	return paramType, optional, defaultValue
}

// NextflowReader reader for workflow written by Nextflow
type NextflowReader struct {
}
