package workflowparser

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils/exec"
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
	if len(validateResultLines) < 3 || !strings.Contains(validateResultLines[len(validateResultLines)-2], "is valid CWL") {
		return "", proto.ErrorWorkflowValidateError("fail to validate workflow")
	}

	depsResults, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--print-deps", path.Join(baseDir, mainWorkflowPath))
	if err != nil {
		return "", err
	}
	// depsResultLines := strings.Split(string(depsResults), "\n")
	// if len(depsResultLines) < 3 {
	// 	return "", proto.ErrorWorkflowValidateError("fail to check deps")
	// }
	// jsonContent := strings.Join(depsResultLines[2:], "\n")
	// var file CWLFile
	// err = json.Unmarshal([]byte(jsonContent), &file)
	// if err != nil {
	// 	return "", err
	// }
	depsFileRe := regexp.MustCompile(`"location":\s*"([^"]*)"`) // 结构比较复杂，没必要解析，直接匹配字符即可
	depsFiles := depsFileRe.FindAllStringSubmatch(string(depsResults), -1)

	workflowFiles := []string{}
	if len(depsFiles) > 0 {
		scriptDir := filepath.Dir(mainWorkflowPath)
		for _, value := range depsFiles {
			joinedPath := filepath.Join(scriptDir, value[1])
			cleanedPath := filepath.Clean(joinedPath)
			if len(cleanedPath) == 0 {
				continue
			}
			// validate file
			if _, err := os.Stat(filepath.Join(baseDir, cleanedPath)); err == nil {
				applog.Infow("file path", "baseDir", baseDir, "relPath", cleanedPath)
				if !sliceContains(workflowFiles, cleanedPath) { // 去掉重复的文件
					workflowFiles = append(workflowFiles, cleanedPath)
				}
			}
		}
	} else {
		return "", apperrors.NewInternalError(fmt.Errorf("no valid workflow files"))
	}

	params := make([]FileParam, 0)
	for _, relPath := range workflowFiles {
		input, err := os.ReadFile(filepath.Join(baseDir, relPath))
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
	rdfResult, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--print-rdf", WorkflowFilePath)
	if err != nil {
		return "", err
	}

	inputsRe := regexp.MustCompile(`(?s)cwl:inputs\s(.*?)\s;`) // 找到所有输入文件
	inputs := inputsRe.FindAllStringSubmatch(string(rdfResult), -1)
	var inputFiles []string
	if len(inputs) > 0 {
		// for _, value := range inputs {
		// 	files := strings.Split(value[1], ",")
		// 	for _, file := range files {
		// 		trimmedFile := strings.TrimSpace(file)
		// 		inputFiles = append(inputFiles, trimmedFile)
		// 	}
		// }
		files := strings.Split(inputs[0][1], ",") // 只处理第一部分，不需要中间文件
		for _, file := range files {
			trimmedFile := strings.TrimSpace(file)
			inputFiles = append(inputFiles, trimmedFile)
		}
	} else {
		return "", nil
	}

	sectionRe := regexp.MustCompile(`[\s]*\r?\n[\s]*\r?\n[\s]*`) // 输出结果分块
	rdfResultSection := sectionRe.Split(string(rdfResult), -1)

	for _, value := range rdfResultSection { // 遍历各个分块，找到描述对应文件的部分
		if findFileSection(value, inputFiles) && !strings.Contains(value, "cwl:inputs") {
			param := WorkflowParam{}
			nameRe := regexp.MustCompile(`#(.*?)> `)
			nameMatch := nameRe.FindStringSubmatch(value)
			if len(nameMatch) > 1 {
				param.Name = nameMatch[1]
				lastSlashIndex := strings.LastIndex(param.Name, "/") // 某些情况下会多一级斜杠
				if lastSlashIndex != -1 {
					param.Name = param.Name[lastSlashIndex+1:]
				}
			} else {
				continue
			}

			param.Optional = strings.Contains(value, "sld:null")

			typeRe := regexp.MustCompile(`sld:type\s+.*?([^:\s,.]+)[\s,.]`)
			typeMatch := typeRe.FindAllStringSubmatch(value, -1)
			if len(typeMatch) > 0 {
				if strings.Contains(typeMatch[0][1], "[") { // Array Enum Record等类型，文本格式与其他类型不同
					param.Type = upperCaseFirstLetter(typeMatch[len(typeMatch)-1][1])
				} else {
					param.Type = upperCaseFirstLetter(typeMatch[0][1])
				}
			} else {
				continue
			}

			defaultRe := regexp.MustCompile(`cwl:default\s+.*?([^;]+)\s`)
			defaultMatch := defaultRe.FindStringSubmatch(value)
			if len(defaultMatch) > 1 && param.Type != "Record" { // Record类型暂时无法输出默认值
				if param.Type == "Array" {
					fmt.Println("[" + removeSpacesAndNewlines(defaultMatch[1]) + "]")
				} else {
					param.Default = defaultMatch[1]
				}
			}
			params = append(params, param)
		}
	}

	return workflowParamDoToPO(params)
}

func (cwl *CWLParser) GetWorkflowOutputs(ctx context.Context, WorkflowFilePath string) (string, error) {
	params := make([]WorkflowParam, 0)
	rdfResult, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--print-rdf", WorkflowFilePath)
	if err != nil {
		return "", err
	}

	outputsRe := regexp.MustCompile(`(?s)cwl:outputs\s(.*?)\s;`) // 找到所有输出文件
	outputs := outputsRe.FindAllStringSubmatch(string(rdfResult), -1)
	var outputFiles []string
	if len(outputs) > 0 {
		// for _, value := range outputs {
		// 	files := strings.Split(value[1], ",")
		// 	for _, file := range files {
		// 		trimmedFile := strings.TrimSpace(file)
		// 		outputFiles = append(outputFiles, trimmedFile)
		// 	}
		// }
		files := strings.Split(outputs[0][1], ",") // 只处理第一部分，不需要中间文件
		for _, file := range files {
			trimmedFile := strings.TrimSpace(file)
			outputFiles = append(outputFiles, trimmedFile)
		}
	} else {
		return "", nil
	}

	sectionRe := regexp.MustCompile(`[\s]*\r?\n[\s]*\r?\n[\s]*`) // 输出结果分块
	rdfResultSection := sectionRe.Split(string(rdfResult), -1)

	for _, value := range rdfResultSection { // 遍历各个分块，找到描述对应文件的部分
		if findFileSection(value, outputFiles) && !strings.Contains(value, "cwl:outputs") {
			param := WorkflowParam{}
			nameRe := regexp.MustCompile(`#(.*?)> `)
			nameMatch := nameRe.FindStringSubmatch(value)
			if len(nameMatch) > 1 {
				param.Name = nameMatch[1]
				lastSlashIndex := strings.LastIndex(param.Name, "/") // 某些情况下会多一级斜杠
				if lastSlashIndex != -1 {
					param.Name = param.Name[lastSlashIndex+1:]
				}
			} else {
				continue
			}

			param.Optional = true // 输出全部为Optional

			typeRe := regexp.MustCompile(`sld:type\s+.*?([^:\s,.]+)[\s,.]`)
			typeMatch := typeRe.FindAllStringSubmatch(value, -1)
			if len(typeMatch) > 0 {
				if strings.Contains(typeMatch[0][1], "[") { // Array Enum Record等类型，文本格式与其他类型不同
					param.Type = upperCaseFirstLetter(typeMatch[len(typeMatch)-1][1])
				} else {
					param.Type = upperCaseFirstLetter(typeMatch[0][1])
				}
			} else {
				continue
			}

			// 不处理default

			params = append(params, param)
		}
	}

	return workflowParamDoToPO(params)
}

func (cwl *CWLParser) GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error) {
	graph, err := exec.Exec(ctx, CommandExecuteTimeout, cwl.Config.CwltoolCmd, "--print-dot", WorkflowFilePath)
	if err != nil { // 无法绘图时不返回错误，使流程能成功导入
		return noGraph, nil
	}

	graphRe := regexp.MustCompile(`(?s)(digraph.*\})`)
	graphMatch := graphRe.FindStringSubmatch(string(graph))

	if len(graphMatch) > 1 {
		return graphMatch[1], nil
	} else {
		return noGraph, nil
	}
}

const noGraph = `digraph G {
    "error" [shape=box, style=filled, color=lightgrey, label="The graph for this workflow is unavailable."];
}`

type CWLFile struct {
	Class          string    `json:"class"`
	Location       string    `json:"location"`
	Format         string    `json:"format"`
	SecondaryFiles []CWLFile `json:"secondaryFiles,omitempty"` // 使用相同的结构体表示嵌套的secondaryFiles数组
	Basename       string    `json:"basename,omitempty"`       // omitempty保证如果字段为空，则在JSON中省略
	Nameroot       string    `json:"nameroot,omitempty"`
	Nameext        string    `json:"nameext,omitempty"`
}

func findFileSection(str string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.HasPrefix(str, substr) {
			return true
		}
	}
	return false
}

func upperCaseFirstLetter(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

func removeSpacesAndNewlines(input string) string {
	parts := strings.Split(input, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return strings.Join(parts, ",")
}

func sliceContains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
