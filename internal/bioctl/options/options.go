package options

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/formatter"
	"github.com/Bio-OS/bioos/internal/bioctl/utils/prompt"
	"github.com/Bio-OS/bioos/pkg/client"
)

const (
	ConfigFlagName    = "config"
	DefaultConfigName = "bioctl.yaml"
)

var DefaultConfig = GlobalOptions{
	Client: ClientOptions{
		Insecure: true,
		Method:   client.HTTPMethod,
		Timeout:  30,
	},
	Stream: StreamOptions{
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	},
}

var CfgFile string

func init() {
	pflag.StringVarP(&CfgFile, ConfigFlagName, "c", CfgFile, "Read bioctl configuration from specified `FILE`, "+
		"support JSON, YAML formats.")
}

type ClientOptions = client.Options

type GlobalOptions struct {
	Client ClientOptions `json:"client" mapstructure:"client"`
	Stream StreamOptions `json:"stream" mapstructure:"stream"`
}

func (o GlobalOptions) Validate() error {
	if err := o.Client.Validate(); err != nil {
		return err
	}
	if err := o.Stream.Validate(); err != nil {
		return err
	}
	return nil
}

type StreamOptions struct {
	// input stream, eg: os.Stdin
	Input io.Reader `json:"input,omitempty"`
	// output stream, eg: os.Stdout
	Output io.Writer `json:"output,omitempty"`
	// error stream, eg: os.Stderr
	ErrorOutput io.Writer `json:"errorOutput,omitempty"`
	// output format, eg: text,table,json
	OutputFormat formatter.Format `json:"format,omitempty" mapstructure:"format,omitempty"`
}

func (o StreamOptions) Validate() error {
	if err := o.OutputFormat.Validate(); err != nil {
		return err
	}
	return nil
}

// LoadConfigFile load config file
func LoadConfigFile(configFile string, opt *GlobalOptions) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(DefaultConfigName)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("conf")
	viper.AddConfigPath(filepath.Join(homedir.HomeDir(), ".bioctl"))
	// Use config file from the flag.
	viper.SetConfigType("yaml") // set the type of the configuration to yaml.

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		utils.CheckErr(err)
	}

	//stream.key should only come from command line, thus delete it here
	if viper.Get("stream") != nil {
		delete(viper.Get("stream").(map[string]interface{}), "format")
	}

	if err := viper.Unmarshal(&opt); err != nil {
		utils.CheckErr(err)
	}
	return
}

func SetConfigEnv() {
	viper.SetEnvPrefix("BIO") // set ENVIRONMENT variables prefix to BIO.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}

type Options interface {
	Complete() error
	Validate() error
	Run(args []string) error
	GetPromptArgs() ([]string, error)
	GetPromptOptions() error
	GetDefaultFormat() formatter.Format
}

func GetCommonRunFunc(o Options) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		utils.CheckErr(o.Complete())

		var err error
		if len(args) > 0 {
			if args[0] == prompt.NeedPromptFlag {
				args, err = o.GetPromptArgs()
				utils.CheckErr(err)
				utils.CheckErr(o.GetPromptOptions())
			}
		}

		utils.CheckErr(o.Validate())
		utils.CheckErr(o.Run(args))
	}
}
