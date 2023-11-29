package workflow

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/schema"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/exec"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type NextflowHandler interface {
	parseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error)
	validateWorkflowFiles(ctx context.Context, version *WorkflowVersion, baseDir, mainWorkflowPath string) error
	getWorkflowInputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error)
	getWorkflowOutputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error)
	getWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error)
}

type Nextflow struct {
	inputParams  []WorkflowParam
	outputParams []WorkflowParam
}

func FileExists(filepath string) bool {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func (next *Nextflow) parseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error) {
	versionRegexp := regexp.MustCompile(VersionRegexpStrNextflow)
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

func (next *Nextflow) validateWorkflowFiles(ctx context.Context, version *WorkflowVersion, baseDir, mainWorkflowPath string) error {
	applog.Infow("start to validate Nextflow files", "mainWorkflowPath", mainWorkflowPath)
	if !FileExists(path.Join(baseDir, mainWorkflowPath)) {
		applog.Errorw("the main workflow file not found", "filepath", path.Join(baseDir, mainWorkflowPath))
		return apperrors.NewInternalError(fmt.Errorf("the main workflow file not found"))
	}
	// mean to `cd $baseDir && nextflow run $mainWorkflowPath -preview`
	result, err := exec.Exec(ctx, CommandExecuteTimeout, baseDir, "nextflow", "run", mainWorkflowPath, "-preview")
	if err != nil {
		applog.Errorw("fail to run with -preview option", "err", err, "result", string(result))
		return apperrors.NewInternalError(fmt.Errorf("run with -preview option failed"))
	}

	// check whether the nextflow_schema.json exists
	if !FileExists(path.Join(baseDir, "nextflow_schema.json")) {
		applog.Errorw("nextflow_schema.json not found", "filepath", path.Join(baseDir, "nextflow_schema.json"))
		return apperrors.NewInternalError(fmt.Errorf("nextflow_schema.json not found"))
	}
	// check whether the nextflow.config exists
	if !FileExists(path.Join(baseDir, "nextflow.config")) {
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

	for process, definition := range nextflowSchema.Definitions {
		for paramName, property := range definition.Properties {
			optional := !utils.In(paramName, definition.Required)
			workflowParam := WorkflowParam{
				Name:     fmt.Sprintf("%s.%s", process, paramName),
				Type:     cases.Title(language.English).String(property.Type),
				Optional: !utils.In(paramName, definition.Required),
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
			applog.Infow("x", "name", paramName, "out", property.Out)
			if !property.Out {
				next.inputParams = append(next.inputParams, workflowParam)
			} else {
				next.outputParams = append(next.outputParams, workflowParam)
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

func (next *Nextflow) getWorkflowInputs(_ context.Context, _ string) ([]WorkflowParam, error) {
	return next.inputParams, nil
}

func (next *Nextflow) getWorkflowOutputs(_ context.Context, _ string) ([]WorkflowParam, error) {
	return next.outputParams, nil
}

func (next *Nextflow) getWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error) {
	//todo get graph
	return "", nil
}
