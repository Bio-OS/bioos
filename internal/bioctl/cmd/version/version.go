package version

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/Bio-OS/bioos/internal/bioctl/factory"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/pkg/version"
)

type VersionOptions struct {
	formatter     formatter.Formatter
	versionClient factory.VersionClient

	option *clioptions.GlobalOptions
}

// NewVersionOptions returns a reference to a CreateOptions
func NewVersionOptions(opt *clioptions.GlobalOptions) *VersionOptions {
	return &VersionOptions{
		option: opt,
	}
}

func NewCmdVersion(opt *clioptions.GlobalOptions) *cobra.Command {
	o := NewVersionOptions(opt)
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version command",
		Long:  `version command`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(o.Complete())
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(args))
		},
	}

	return cmd
}

func (o *VersionOptions) Complete() error {
	var err error
	f := factory.NewFactory(&o.option.Client)
	o.versionClient, err = f.VersionClient()
	if err != nil {
		return err
	}
	o.formatter = formatter.NewFormatter(o.option.Stream.OutputFormat, o.option.Stream.Output)
	return nil
}

func (o *VersionOptions) Validate() error {
	//return o.option.Validate()
	return nil
}

func (o *VersionOptions) Run(args []string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(o.option.Client.Timeout)*time.Second)
	defer cancelFunc()

	switch o.formatter.(type) {
	case *formatter.Json, *formatter.Yaml:
		m := map[string]interface{}{
			"clientVersion": version.Get(),
		}
		info, err := o.versionClient.Version(ctx)
		if err != nil {
			m["serverVersion"] = err.Error()
		} else {
			m["serverVersion"] = info
		}
		o.formatter.Write(m)

	case *formatter.Table:
		// todo optimize here
		o.formatter.Write(version.Get())
		info, err := o.versionClient.Version(ctx)
		if err != nil {
			o.formatter.Write(err.Error())
		} else {
			o.formatter.Write(info)
		}

	case *formatter.Text:
		o.formatter.Write("clientVersion")
		o.formatter.Write(version.Get())
		o.formatter.Write("serverVersion")
		info, err := o.versionClient.Version(ctx)
		if err != nil {
			o.formatter.Write(err.Error())
		} else {
			o.formatter.Write(info)
		}
	}
	return nil
}
