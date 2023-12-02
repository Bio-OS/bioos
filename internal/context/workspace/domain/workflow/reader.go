package workflow

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/schema"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/exec"
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

var (
	_ Reader = &WDLReader{}
	_ Reader = &NextflowReader{}
)

// WDLReader reader for workflow written by WDL
type WDLReader struct {
	options *ReaderOptions
}

// ParseWorkflowVersion ...
func (r *WDLReader) ParseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error) {
	versionRegexp := regexp.MustCompile(VersionWDLRegexpStr)
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
	validateResult, err := exec.Exec(ctx, CommandExecuteTimeout, ".", "java", "-jar", r.options.WomtoolPath, "validate", path.Join(baseDir, mainWorkflowPath), "-l")
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
	graph, err := exec.Exec(ctx, CommandExecuteTimeout, ".", "java", "-jar", r.options.WomtoolPath, "graph", WorkflowFilePath)
	if err != nil {
		return "", err
	}

	return string(graph), nil
}

// getWorkflowParams ...
func (r *WDLReader) getWorkflowParams(ctx context.Context, name string, arg ...string) ([]WorkflowParam, error) {
	params := make([]WorkflowParam, 0)
	outputsResult, err := exec.Exec(ctx, CommandExecuteTimeout, ".", name, arg...)
	if err != nil {
		return params, err
	}
	var outputsMap map[string]string
	if err := json.Unmarshal(outputsResult, &outputsMap); err != nil {
		return params, err
	}

	for paramName, value := range outputsMap {
		paramType, optional, defaultValue := r.parseWorkflowParamValue(value)
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

func (_ *WDLReader) parseWorkflowParamValue(value string) (paramType string, optional bool, defaultValue *string) {
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
	inputParams   []WorkflowParam
	outputParams  []WorkflowParam
	graphFilepath string
}

// ValidateWorkflowFiles ...
func (r *NextflowReader) ValidateWorkflowFiles(ctx context.Context, version *WorkflowVersion, baseDir, mainWorkflowPath string) error {
	// check whether the main workflow file exists
	if mainWorkflowPath != "main.nf" {
		applog.Errorw("wrong name", "name", mainWorkflowPath)
		return apperrors.NewInternalError(fmt.Errorf("the name of main workflow file is wrong"))
	}
	if !utils.FileExists(path.Join(baseDir, mainWorkflowPath)) {
		applog.Errorw("the main workflow file not found", "filepath", path.Join(baseDir, mainWorkflowPath))
		return apperrors.NewInternalError(fmt.Errorf("the main workflow file not found"))
	}
	// mean to `cd $baseDir && nextflow run $mainWorkflowPath -preview`
	result, err := exec.Exec(ctx, CommandExecuteTimeout, baseDir, "nextflow", "run", mainWorkflowPath, "-preview", "-with-dag", "dag.html")
	if err != nil {
		applog.Errorw("fail to run with -preview option", "err", err, "result", string(result))
		return apperrors.NewInternalError(fmt.Errorf("run with -preview option failed"))
	}
	r.graphFilepath = path.Join(baseDir, "dag.html")

	// check whether the nextflow_schema.json exists
	if !utils.FileExists(path.Join(baseDir, "nextflow_schema.json")) {
		applog.Errorw("nextflow_schema.json not found", "filepath", path.Join(baseDir, "nextflow_schema.json"))
		return apperrors.NewInternalError(fmt.Errorf("nextflow_schema.json not found"))
	}
	// check whether the nextflow.config exists
	if !utils.FileExists(path.Join(baseDir, "nextflow.config")) {
		applog.Errorw("nextflow.config not found", "filepath", path.Join(baseDir, "nextflow.config"))
		return apperrors.NewInternalError(fmt.Errorf("nextflow.config not found"))
	}
	// try to parse the nextflow_schema.json
	f, err := os.Open(path.Join(baseDir, "nextflow_schema.json"))
	if err != nil {
		applog.Errorw("fail to open the schema file", "err", err)
		return apperrors.NewInternalError(fmt.Errorf("nextflow_schema.json not found"))
	}

	nextflowSchema := schema.NextflowSchema{}
	err = json.NewDecoder(f).Decode(&nextflowSchema)
	if err != nil {
		applog.Errorw("fail to unmarshal the input json", "err", err)
		return apperrors.NewInternalError(fmt.Errorf("please check the format of nextflow_schema.json"))
	}

	for _, definition := range nextflowSchema.Definitions {
		for paramName, property := range definition.Properties {
			optional := !utils.In(paramName, definition.Required)
			workflowParam := WorkflowParam{
				Name:     paramName,
				Type:     cases.Title(language.English).String(property.Type),
				Optional: optional,
			}
			if property.Mimetype != "" {
				workflowParam.Type = "File"
			}
			if !optional && property.Default != nil {
				switch property.Type {
				case schema.TypeString:
					workflowParam.Default = property.Default.(string)
				case schema.TypeBoolean:
					workflowParam.Default = strconv.FormatBool(property.Default.(bool))
				case schema.TypeNumber:
					workflowParam.Default = strconv.FormatFloat(property.Default.(float64), 64, 6, 10)
				case schema.TypeInteger:
					workflowParam.Default = strconv.FormatInt(property.Default.(int64), 10)
				}
			}

			if !property.Out {
				r.inputParams = append(r.inputParams, workflowParam)
			} else {
				r.outputParams = append(r.outputParams, workflowParam)
			}
		}
	}

	// add workflow file
	workflowFiles := []string{"nextflow_schema.json", "nextflow.config"}
	if err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".nf") {
			// TODO need to support multi-level directory
			workflowFiles = append(workflowFiles, path[len(baseDir)+1:])
		}
		return nil
	}); err != nil {
		applog.Errorw("fail to lookup .nf files", "baseDir", baseDir)
		return apperrors.NewInternalError(err)
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
func (r *NextflowReader) GetWorkflowInputs(_ context.Context, _ string) ([]WorkflowParam, error) {
	return r.inputParams, nil
}

// GetWorkflowOutputs ...
func (r *NextflowReader) GetWorkflowOutputs(_ context.Context, _ string) ([]WorkflowParam, error) {
	return r.outputParams, nil
}

// GetWorkflowGraph ...
func (r *NextflowReader) GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error) {
	htmlFile, err := os.Open(r.graphFilepath)
	if err != nil {
		applog.Errorw("fail to open dag.html", "err", err)
		return "", apperrors.NewInternalError(fmt.Errorf("open dag.html failed"))
	}
	defer htmlFile.Close()
	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		applog.Errorw("fail to parse dag.htm", "err", err)
		return "", apperrors.NewInternalError(fmt.Errorf("parse dag.html failed"))
	}
	return doc.Find(".mermaid").Eq(0).Text(), nil
}

// ParseWorkflowVersion ...
func (r *NextflowReader) ParseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error) {
	versionRegexp := regexp.MustCompile(VersionNextflowRegexpStr)
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
		fmt.Println(line, matched)
		if matched != nil && len(matched) >= 2 {
			applog.Infow("version regexp matched", "matched", matched)
			return "DSL" + matched[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "DSL2", nil
}
