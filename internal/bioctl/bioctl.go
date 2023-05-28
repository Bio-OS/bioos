package bioctl

import (
	"flag"
	"io"
	"os"
	"reflect"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	cliflag "k8s.io/component-base/cli/flag"

	internalcmd "github.com/Bio-OS/bioos/internal/bioctl/cmd"
	clidatamodel "github.com/Bio-OS/bioos/internal/bioctl/cmd/data-model"
	clisubmission "github.com/Bio-OS/bioos/internal/bioctl/cmd/submission"
	cliversion "github.com/Bio-OS/bioos/internal/bioctl/cmd/version"
	cliworkflow "github.com/Bio-OS/bioos/internal/bioctl/cmd/workflow"
	cliworkspace "github.com/Bio-OS/bioos/internal/bioctl/cmd/workspace"
	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
)

// NewDefaultBIOCtlCommand return the bioctl cmd.
func NewDefaultBIOCtlCommand() *cobra.Command {
	return NewBIOCtlCommandWithArgs(os.Stdin, os.Stdout, os.Stderr)
}

// NewBIOCtlCommandWithArgs initialized a bioctl command.
func NewBIOCtlCommandWithArgs(in io.Reader, out, err io.Writer) *cobra.Command {
	command := newBIOCtlCommand()
	command.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	flags := command.PersistentFlags()

	addProfilingFlags(flags)
	command.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

	opt := clioptions.DefaultConfig

	command.AddCommand(cliworkspace.NewCmdWorkspace(&opt))
	command.AddCommand(cliworkflow.NewCmdWorkflow(&opt))
	command.AddCommand(clidatamodel.NewCmdDataModel(&opt))
	command.AddCommand(clisubmission.NewCmdSubmission(&opt))
	addExample(command)

	// version doesn't need Example text
	command.AddCommand(cliversion.NewCmdVersion(&opt))

	command.PersistentFlags().AddFlag(pflag.Lookup(clioptions.ConfigFlagName))
	// read from config file
	clioptions.LoadConfigFile(clioptions.CfgFile, &opt)

	clioptions.SetConfigEnv()

	// define global options
	loadConfigFromOptionOrEnv(command, &opt.Client.ServerAddr, "server-addr", "SERVER_ADDR", "Bioos apiserver address")
	loadConfigFromOptionOrEnv(command, &opt.Client.Insecure, "insecure", "INSECURE", "Whether enable tls")
	loadConfigFromOptionOrEnv(command, &opt.Client.ServerCertFile, "server-cert-file", "SERVER_CERT_FILE", "Server cert file path")
	loadConfigFromOptionOrEnv(command, &opt.Client.ServerName, "server-name", "SERVER_NAME", "Server Name")
	loadConfigFromOptionOrEnv(command, &opt.Client.ClientCertFile, "client-cert-file", "CLIENT_CERT_FILE", "Client cert file path")
	loadConfigFromOptionOrEnv(command, &opt.Client.ClientCertKeyFile, "client-cert-key-file", "CLIENT_CERT_KEY_FILE", "Client key file path")
	loadConfigFromOptionOrEnv(command, &opt.Client.CaFile, "ca-file", "CA_FILE", "CA file path")
	loadConfigFromOptionOrEnv(command, &opt.Client.Username, "username", "USERNAME", "Username")
	loadConfigFromOptionOrEnv(command, &opt.Client.Password, "password", "PASSWORD", "Password")
	loadConfigFromOptionOrEnv(command, &opt.Client.AuthToken, "auth-token", "AUTH_TOKEN", "Auth token")
	loadConfigFromOptionOrEnv(command, (*string)(&opt.Client.Method), "connect-method", "CONNECT_METHOD", "Cli connect method: grpc or http")
	loadConfigFromOptionOrEnv(command, &opt.Client.Timeout, "connect-timeout", "CONNECT_TIMEOUT", "Connect timeout seconds")

	// format should only come from command line
	command.PersistentFlags().StringVarP((*string)(&opt.Stream.OutputFormat), "output-format", "o", string(opt.Stream.OutputFormat), "Cli output format: json, yaml, table or text")

	_ = viper.BindPFlags(command.PersistentFlags())
	// read from environment

	command.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return command
}

var castFuncMap = map[reflect.Kind]interface{}{
	reflect.Bool:   cast.ToBool,
	reflect.String: cast.ToString,
	reflect.Int:    cast.ToInt,
}

var pfFuncMap = map[reflect.Kind]interface{}{
	reflect.Bool:   (*pflag.FlagSet).BoolVar,
	reflect.String: (*pflag.FlagSet).StringVar,
	reflect.Int:    (*pflag.FlagSet).IntVar,
}

// Read config from env or flag, flag > env
func loadConfigFromOptionOrEnv(command *cobra.Command, opt interface{}, optName, envName, usage string) {
	optT := reflect.TypeOf(opt)
	for optT.Kind() == reflect.Pointer {
		optT = optT.Elem()
	}
	optKind := optT.Kind()
	viper.BindEnv(envName)
	eVal := viper.Get(envName)
	defaultVal := reflect.ValueOf(opt).Elem()
	if eVal != nil {
		defaultVal = reflect.ValueOf(castFuncMap[optKind]).Call([]reflect.Value{reflect.ValueOf(eVal)})[0]
	}
	reflect.ValueOf(pfFuncMap[optKind]).Call([]reflect.Value{
		reflect.ValueOf(command.PersistentFlags()),
		reflect.ValueOf(opt),
		reflect.ValueOf(optName),
		defaultVal,
		reflect.ValueOf(usage),
	},
	)
}

func newBIOCtlCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "bioctl",
		Short: "bioctl command",
		Long:  `bioctl command`,
		Run:   runHelp,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return flushProfiling()
		},
	}

}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func addExample(cmd *cobra.Command) {
	if !cmd.HasSubCommands() {
		if cmd.Example != "" {
			return
		}
		cmd.Example = internalcmd.GenExample(cmd)
		return
	}
	for _, subCommand := range cmd.Commands() {
		addExample(subCommand)
	}
}
