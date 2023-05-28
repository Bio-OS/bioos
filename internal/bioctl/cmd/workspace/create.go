package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
)

// CreateOptions is an options to create a workspace.
type CreateOptions struct {
	Description string
	MountType   string
	MountPath   string

	workspaceClient factory.WorkspaceClient
	formatter       formatter.Formatter

	options *clioptions.GlobalOptions
}

// NewCreateOptions returns a reference to a CreateOptions
func NewCreateOptions(opt *clioptions.GlobalOptions) *CreateOptions {
	return &CreateOptions{
		options: opt,
	}
}

func NewCmdCreate(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewCreateOptions(opt)

	cmd := &cobra.Command{
		Use:   "create <workspace_name>",
		Short: "create a workspace",
		Long:  "create a workspace",
		Args:  cobra.ExactArgs(1),
		Run:   clioptions.GetCommonRunFunc(o),
	}

	cmd.Flags().StringVarP(&o.Description, "description", "d", o.Description, "The description of the workspace.")
	cmd.Flags().StringVarP(&o.MountType, "mount-type", "t", o.MountType, "The mount type of the workspace Storage.")
	cmd.Flags().StringVarP(&o.MountPath, "mount-path", "p", o.MountPath, "The mount path of the workspace Storage.")

	return cmd
}

// Complete completes all the required options.
func (o *CreateOptions) Complete() error {
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

// Validate validate the create options
func (o *CreateOptions) Validate() error {
	if err := o.options.Validate(); err != nil {
		return err
	}
	if o.MountType != "nfs" {
		return fmt.Errorf("workspace storage [%s] not support", o.MountType)
	}
	return nil
}

// Run run the create workspace command
func (o *CreateOptions) Run(args []string) error {
	workspaceName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.Client.Timeout))
	defer cancel()
	req := &convert.CreateWorkspaceRequest{
		Name:        workspaceName,
		Description: o.Description,
	}
	if o.MountPath != "" {
		req.Storage = &convert.WorkspaceStorage{
			NFS: &convert.NFSWorkspaceStorage{
				MountPath: o.MountPath,
			},
		}
	}

	resp, err := o.workspaceClient.CreateWorkspace(ctx, req)
	if err != nil {
		return err
	}
	o.formatter.Write(resp.Id)

	return nil
}

func (o *CreateOptions) GetPromptArgs() ([]string, error) {
	workspaceName, err := prompt.PromptRequiredString("Name")
	if err != nil {
		return []string{}, err
	}
	return []string{workspaceName}, nil
}

func (o *CreateOptions) GetPromptOptions() error {
	var err error
	o.Description, err = prompt.PromptRequiredString("Description")
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

func (o *CreateOptions) GetDefaultFormat() formatter.Format {
	return formatter.TextFormat
}
