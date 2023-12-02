package workflowparser

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils/exec"
	"gopkg.in/yaml.v2"
)

type CWLConfig struct {
	CwltoolCmd string
}

type CWLParser struct {
	Config CWLConfig
}

func NewCWLParser(config CWLConfig) *CWLParser {
	return &CWLParser{Config: config}
}

func (cwl *CWLParser) ParseWorkflowVersion(_ context.Context, mainWorkflowPath string) (string, error) {
	versionRegexp := regexp.MustCompile(CWLVersionRegexpStr)
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
		if len(matched) >= 2 {
			applog.Infow("version regexp matched", "matched", matched)
			return matched[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "v1.0", nil
}

func (cwl *CWLParser) ValidateWorkflowFiles(ctx context.Context, baseDir, mainWorkflowPath string) (string, error) {
	applog.Infow("start to validate files", "mainWorkflowPath", mainWorkflowPath)
	validateResult, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--validate", path.Join(baseDir, mainWorkflowPath))
	if err != nil {
		applog.Errorw("fail to validate workflow", "err", err, "result", string(validateResult))
		return "", apperrors.NewInternalError(fmt.Errorf("validate workflow failed"))
	}
	validateResultLines := strings.Split(string(validateResult), "\n")
	applog.Infow("validate result", "result", validateResultLines)
	// parse and save workflow files
	if len(validateResultLines) < 3 || !strings.Contains(validateResultLines[2], "is valid CWL") {
		return "", proto.ErrorWorkflowValidateError("fail to validate workflow")
	}

	depsResults, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--print-deps", path.Join(baseDir, mainWorkflowPath))
	if err != nil {
		return "", err
	}
	depsResultLines := strings.Split(string(depsResults), "\n")
	if len(validateResultLines) < 3 {
		return "", proto.ErrorWorkflowValidateError("fail to check deps")
	}
	jsonContent := strings.Join(depsResultLines[2:], "\n")
	var file CWLFile
	err = json.Unmarshal([]byte(jsonContent), &file)
	if err != nil {
		return "", err
	}

	workflowFiles := []string{mainWorkflowPath}
	// need to start from line 2(start with line 0)
	for i := 0; i < len(file.SecondaryFiles); i++ {
		workflowFiles = append(workflowFiles, file.SecondaryFiles[i].Location)
	}
	params := make([]FileParam, 0)
	for _, relPath := range workflowFiles {
		input, err := os.ReadFile(path.Join(baseDir, relPath))
		if err != nil {
			applog.Errorw("fail to read file", "err", err)
			return "", apperrors.NewInternalError(err)
		}

		encodedContent := base64.StdEncoding.EncodeToString(input)

		param := FileParam{
			Path:    relPath,
			Content: encodedContent,
		}
		params = append(params, param)
	}
	return fileParamDoToPO(params)
}

func (cwl *CWLParser) GetWorkflowInputs(ctx context.Context, WorkflowFilePath string) (string, error) {
	params := make([]WorkflowParam, 0)
	templateResult, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--make-template", WorkflowFilePath)
	if err != nil {
		return "", err
	}
	templateResultLines := strings.Split(string(templateResult), "\n")
	yamlContent := strings.Join(templateResultLines[2:], "\n")

	var yamlData map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &yamlData); err != nil {
		return "", err
	}

	for key, value := range yamlData {
		param := WorkflowParam{Name: key}

		keyWithColon := key + ":"
		startIndex := strings.Index(yamlContent, keyWithColon)
		if startIndex == -1 {
			continue
		}
		endIndex := strings.Index(yamlContent[startIndex:], "\n")
		var paramDesc string
		if endIndex == -1 {
			paramDesc = yamlContent[startIndex:]
		} else {
			paramDesc = yamlContent[startIndex:endIndex]
		}

		reType := regexp.MustCompile(`type\s+'([^']*)'`)
		matchType := reType.FindStringSubmatch(paramDesc)
		if matchType != nil {
			param.Type = matchType[1]
		} else {
			continue
		}

		param.Optional = strings.Contains(paramDesc, "(optional)")

		if strings.Contains(paramDesc, "default value") {
			param.Default = value.(string)
		}

		params = append(params, param)
	}

	return workflowParamDoToPO(params)
}

func (cwl *CWLParser) GetWorkflowOutputs(ctx context.Context, WorkflowFilePath string) (string, error) {
	params := make([]WorkflowParam, 0)
	rdfResult, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--print-rdf", WorkflowFilePath)
	if err != nil {
		return "", err
	}
	rdfResultSection := strings.Split(string(rdfResult), "\n\n")
	for _, value := range rdfResultSection {
		if strings.Contains(value, "cwl:outputBinding") {
			param := WorkflowParam{Optional: false}
			re := regexp.MustCompile(`sld:type\s+([^:\s]+):`)
			match := re.FindStringSubmatch(value)
			if len(match) > 1 {
				param.Type = match[1]
			} else {
				continue
			}

			re2 := regexp.MustCompile(`\/(.*?)> rdfs:comment`)
			match2 := re2.FindStringSubmatch(value)
			if len(match2) > 1 {
				param.Name = match[1]
			} else {
				continue
			}
			param.Default = value
			params = append(params, param)
		}
	}

	return workflowParamDoToPO(params)
}

func (cwl *CWLParser) GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error) {
	return "", nil
}

type CWLFile struct {
	Class          string    `json:"class"`
	Location       string    `json:"location"`
	Format         string    `json:"format"`
	SecondaryFiles []CWLFile `json:"secondaryFiles,omitempty"` // 使用相同的结构体表示嵌套的secondaryFiles数组
	Basename       string    `json:"basename,omitempty"`       // omitempty保证如果字段为空，则在JSON中省略
	Nameroot       string    `json:"nameroot,omitempty"`
	Nameext        string    `json:"nameext,omitempty"`
}
