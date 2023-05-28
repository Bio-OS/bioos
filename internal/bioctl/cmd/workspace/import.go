package workspace

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/pkg/consts"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/schema"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type ImportOptions struct {
	YamlPath  string
	MountType string
	MountPath string

	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewImportOptions returns a reference to a ImportOptions
func NewImportOptions(opt *clioptions.GlobalOptions) *ImportOptions {
	return &ImportOptions{
		options: opt,
	}
}

func NewCmdImport(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewImportOptions(opt)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "import a workspace",
		Long:  "import a workspace",
		Args:  cobra.NoArgs,
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.YamlPath, "yaml", "y", o.YamlPath, "The path of the workspace yaml file")
	cmd.Flags().StringVarP(&o.MountType, "mount-type", "t", o.MountType, "The mount type of the workspace Storage.")
	cmd.Flags().StringVarP(&o.MountPath, "mount-path", "p", o.MountPath, "The mount path of the workspace Storage.")

	return cmd
}

// Complete completes all the required options.
func (o *ImportOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.options.Client)
	o.workspaceClient, err = f.WorkspaceClient()
	if err != nil {
		return err
	}

	if o.options.Stream.OutputFormat == "" {
		o.options.Stream.OutputFormat = o.GetDefaultFormat()
	}
	o.formatter = formatter.NewFormatter(o.options.Stream.OutputFormat, o.options.Stream.Output)
	return nil
}

// Validate validate the import options
func (o *ImportOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.YamlPath == "" || path.Ext(o.YamlPath) != ".yaml" {
		return fmt.Errorf("you must choose a csv file to import")
	}
	_, err := os.Stat(o.YamlPath)
	if err != nil {
		return err
	}
	if o.MountType != "nfs" {
		return fmt.Errorf("workspace storage [%s] not support", o.MountType)
	}
	return nil
}

// Run run the import workspace command
func (o *ImportOptions) Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()

	bytes, err := os.ReadFile(o.YamlPath)
	if err != nil {
		return err
	}
	schema := &schema.WorkspaceTypedSchema{}
	err = yaml.Unmarshal(bytes, &schema)
	if err != nil {
		return err
	}

	zipFilePath, err := o.zipFiles(schema)
	if err != nil {
		return err
	}
	defer os.Remove(zipFilePath)

	req := &convert.ImportWorkspaceRequest{
		FilePath:  zipFilePath,
		MountType: o.MountType,
	}
	if o.MountPath != "" {
		req.MountPath = o.MountPath
	}

	resp, err := o.workspaceClient.ImportWorkspace(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp.Id)

	return nil
}

func (o *ImportOptions) GetPromptArgs() ([]string, error) {
	return nil, nil
}

func (o *ImportOptions) GetPromptOptions() error {
	var err error

	o.YamlPath, err = prompt.PromptStringWithValidator("Workspace YamlPath File Path", func(ans interface{}) error {
		err = survey.Required(ans)
		if err != nil {
			return err
		}
		_, err = os.Stat(cast.ToString(ans))
		if err != nil {
			curPath, _ := os.Getwd()
			return fmt.Errorf("%w, (currenct path is [%s])", err, curPath)
		}
		return nil
	}, prompt.WithInputMessage("only support YamlPath"))

	if err != nil {
		return err
	}

	o.MountPath, err = prompt.PromptRequiredString("MountPath", prompt.WithInputMessage("abs path"))
	if err != nil {
		return err
	}

	o.MountType, err = prompt.PromptStringSelect("MountType", 1, []string{"nfs"})
	if err != nil {
		return err
	}

	return nil
}

func (o *ImportOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}

func (o *ImportOptions) zipFiles(schema *schema.WorkspaceTypedSchema) (string, error) {
	yamlBase := path.Base(o.YamlPath)
	yamlFilePrefix := strings.TrimSuffix(yamlBase, path.Ext(o.YamlPath))
	dir := fmt.Sprintf("import-%s-%s-%s", yamlFilePrefix, o.MountType, rand.String(10))
	os.RemoveAll(dir)

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		applog.Errorf("fail to make tmp dir: %s", err)
		return "", err
	}
	defer os.RemoveAll(dir)

	baseDir := filepath.Dir(o.YamlPath)
	// copy all files to temp dir
	filesMap := map[string]string{o.YamlPath: path.Join(dir, consts.WorkspaceYAMLName)}
	for _, dataModel := range schema.DataModels {
		filesMap[path.Join(baseDir, dataModel.Path)] = path.Join(dir, dataModel.Path)
	}
	for _, workflow := range schema.Workflows {
		err = filepath.Walk(path.Join(baseDir, workflow.Path), func(p string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			filesMap[p] = path.Join(dir, workflow.Path, path.Base(p))
			return nil
		})
		if err != nil {
			return "", err
		}

	}
	for _, artifact := range schema.Notebooks.Artifacts {
		filesMap[path.Join(baseDir, artifact.Path)] = path.Join(dir, artifact.Path)
	}
	for fToRead, fToWrite := range filesMap {
		input, err := os.ReadFile(filepath.Clean(fToRead))
		if err != nil {
			return "", err
		}
		if _, err = os.Stat(fToWrite); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(fToWrite), os.ModePerm)
			if err != nil {
				return "", err
			}
		}
		err = os.WriteFile(fToWrite, input, 0660)
		if err != nil {
			return "", err
		}
	}

	zipFilePath := strings.ReplaceAll(yamlBase, ".yaml", ".zip")
	err = utils.ZipDir(dir, zipFilePath)
	if err != nil {
		return "", err
	}
	return zipFilePath, nil
}
