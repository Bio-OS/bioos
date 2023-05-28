package submission

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/cmd"
	clidatamodel "github.com/Bio-OS/bioos/internal/bioctl/cmd/data-model"
	cliworkflow "github.com/Bio-OS/bioos/internal/bioctl/cmd/workflow"
	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/pkg/consts"
)

type InputsOutputs struct {
	InputsTemplate  map[string]interface{}
	OutputsTemplate map[string]interface{}
	InputsMaterial  map[string]interface{}
	OutputsMaterial map[string]interface{}
}

const emptyJsonStr = "null"

// SubmitOptions is an options to submit a workspace.
type SubmitOptions struct {
	WorkspaceName   string
	Description     string
	Type            string
	DataModelName   string
	DataModelRowIDs []string
	File            string
	ReadFromCache   bool

	InputsTemplate  string
	OutputsTemplate string
	InputsMaterial  string
	OutputsMaterial string

	dataModelClient  factory.DataModelClient
	workflowClient   factory.WorkflowClient
	submissionClient factory.SubmissionClient
	workspaceClient  factory.WorkspaceClient
	formatter        formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewSubmitOptions returns a reference to a SubmitOptions
func NewSubmitOptions(opt *clioptions.GlobalOptions) *SubmitOptions {
	return &SubmitOptions{
		options: opt,
	}
}

func NewCmdSubmit(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewSubmitOptions(opt)

	cmd := &cobra.Command{
		Use:   "submit <workflow_name>",
		Short: "submit a workflow",
		Long:  "submit a workflow",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.WorkspaceName, "workspace", "w", o.WorkspaceName, "The workspace name")
	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "The description of the submission.")
	cmd.Flags().StringVarP(&o.Type, "type", "t", o.Type, "The Type of the submission.")
	cmd.Flags().StringVarP(&o.DataModelName, "data-model", "m", o.DataModelName, "The name of the data-model this submission will use.")
	cmd.Flags().StringSliceVar(&o.DataModelRowIDs, "data-model-rows", o.DataModelRowIDs, "The rows of the data-model this submission will use.")
	cmd.Flags().StringVarP(&o.File, "file", "f", o.File, "The file path of Inputs/Outputs.")
	cmd.Flags().BoolVar(&o.ReadFromCache, "call-caching", true, "use previous cache of the submission or not.")

	return cmd
}

// Complete completes all the required options.
func (o *SubmitOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workspaceClient, err = f.WorkspaceClient()
	if err != nil {
		return err
	}
	o.submissionClient, err = f.SubmissionClient()
	if err != nil {
		return err
	}
	o.dataModelClient, err = f.DataModelClient()
	if err != nil {
		return err
	}
	o.workflowClient, err = f.WorkflowClient()
	if err != nil {
		return err
	}

	if o.options.Stream.OutputFormat == "" {
		o.options.Stream.OutputFormat = o.GetDefaultFormat()
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the submit options
func (o *SubmitOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}

	if o.WorkspaceName == "" {
		return fmt.Errorf("need to specify a workspace name")
	}

	if o.Type == "" {
		return fmt.Errorf("submission type cannot be empty")
	}

	if o.File == "" {
		return fmt.Errorf("need to specify a file to declare inputs and outputs")
	}

	err := o.parseInputsAndOutputsFile()
	if err != nil {
		return err
	}

	if o.Type == consts.DataModelTypeSubmission {
		if o.InputsTemplate == "" {
			return fmt.Errorf("InputsTemplate cannot be empty")
		}
		if o.DataModelName == "" {
			return fmt.Errorf("need to specify a data-model")
		}
	} else if o.Type == consts.FilePathTypeSubmission {
		if o.InputsMaterial == "" {
			return fmt.Errorf("InputsMaterial cannot be empty")
		}
	} else {
		return fmt.Errorf("submission type %s not support", o.Type)
	}

	return nil
}

// Run run the submit workspace command
func (o *SubmitOptions) Run(args []string) error {
	workflowName := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	workspaceID, err := cliworkspace.ConvertWorkspaceNameIntoID(ctx, o.workspaceClient, o.WorkspaceName)
	if err != nil {
		return err
	}

	workflowID, err := cliworkflow.ConvertWorkflowNameIntoID(ctx, o.workflowClient, workspaceID, workflowName)
	if err != nil {
		return err
	}

	req := &convert.CreateSubmissionRequest{
		WorkspaceID: workspaceID,
		Name:        fmt.Sprintf("%s-history-%s", workflowName, time.Now().Format("2006-01-02-15-04-05")),
		WorkflowID:  workflowID,
		Type:        o.Type,
		ExposedOptions: convert.ExposedOptions{
			ReadFromCache: o.ReadFromCache,
		},
	}

	if o.Type == consts.DataModelTypeSubmission {
		dataModelID, err := clidatamodel.ConvertDataModelNameIntoID(ctx, o.dataModelClient, workspaceID, o.DataModelName)
		if err != nil {
			return err
		}
		if len(o.DataModelRowIDs) == 0 {
			idsResp, err := o.dataModelClient.ListAllDataModelRowIDs(ctx, &convert.ListAllDataModelRowIDsRequest{
				WorkspaceID: workspaceID,
				ID:          dataModelID,
			})
			if err != nil {
				return err
			}
			o.DataModelRowIDs = idsResp.RowIDs
		}
		req.Entity = &convert.Entity{
			DataModelID:     dataModelID,
			DataModelRowIDs: o.DataModelRowIDs,
			InputsTemplate:  o.InputsTemplate,
			OutputsTemplate: o.OutputsTemplate,
		}
	} else if o.Type == consts.FilePathTypeSubmission {
		req.InOutMaterial = &convert.InOutMaterial{
			InputsMaterial:  o.InputsMaterial,
			OutputsMaterial: o.OutputsMaterial,
		}
	}

	if o.Description != "" {
		req.Description = &o.Description
	}

	resp, err := o.submissionClient.CreateSubmission(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp.ID)

	return nil
}

func (o *SubmitOptions) GetPromptArgs() ([]string, error) {
	workflowName, err := prompt.PromptRequiredString("Workflow Name")
	if err != nil {
		return []string{}, err
	}

	return []string{workflowName}, nil
}

func (o *SubmitOptions) GetPromptOptions() error {
	var err error
	// TODO put it in GetPromptArgs? enhance interactive user experience
	o.WorkspaceName, err = cmd.GetWorkspaceName(o.options.Client.Timeout, o.workspaceClient)
	if err != nil {
		return err
	}

	o.Description, err = prompt.PromptOptionalString("Description")
	if err != nil {
		return err
	}

	o.Type, err = prompt.PromptStringSelect("Submission Type", 2, []string{consts.DataModelTypeSubmission, consts.FilePathTypeSubmission})
	if err != nil {
		return err
	}

	if o.Type == consts.DataModelTypeSubmission {
		o.DataModelName, err = prompt.PromptRequiredString("DataModel Name")
		if err != nil {
			return err
		}

		o.DataModelRowIDs, err = prompt.PromptStringSlice(fmt.Sprintf("RowIDs of [%s]", o.DataModelName))
		if err != nil {
			return err
		}
	}

	o.File, err = prompt.PromptStringWithValidator("Inputs and Outputs FilePath", func(ans interface{}) error {
		err := survey.Required(ans)
		if err != nil {
			return err
		}
		_, err = os.Stat(cast.ToString(ans))
		if err != nil {
			curPath, _ := os.Getwd()
			return fmt.Errorf("%w, (currenct path is [%s])", err, curPath)
		}
		return nil
	})
	if err != nil {
		return err
	}

	o.ReadFromCache, err = prompt.PromptBoolSelect("ReadFromCache")
	if err != nil {
		return err
	}

	return nil
}

func (o *SubmitOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}

func (o *SubmitOptions) parseInputsAndOutputsFile() error {
	bytes, err := os.ReadFile(o.File)
	if err != nil {
		return err
	}
	inOuts := InputsOutputs{}
	err = json.Unmarshal(bytes, &inOuts)
	if err != nil {
		return err
	}
	dataInputsTemplate, _ := json.Marshal(inOuts.InputsTemplate)
	o.InputsTemplate = string(dataInputsTemplate)
	dataOutputsTemplate, _ := json.Marshal(inOuts.OutputsTemplate)
	o.OutputsTemplate = string(dataOutputsTemplate)
	dataInputsMaterial, _ := json.Marshal(inOuts.InputsMaterial)
	o.InputsMaterial = string(dataInputsMaterial)
	dataOutputsMaterial, _ := json.Marshal(inOuts.OutputsMaterial)
	o.OutputsMaterial = string(dataOutputsMaterial)

	if o.InputsTemplate == emptyJsonStr {
		o.InputsTemplate = ""
	}
	if o.OutputsTemplate == emptyJsonStr {
		o.OutputsTemplate = ""
	}
	if o.InputsMaterial == emptyJsonStr {
		o.InputsMaterial = ""
	}
	if o.OutputsMaterial == emptyJsonStr {
		o.OutputsMaterial = ""
	}

	return nil
}
