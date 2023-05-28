package options

import (
	"fmt"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	"github.com/Bio-OS/bioos/pkg/auth"
	"github.com/Bio-OS/bioos/pkg/client"
	"github.com/Bio-OS/bioos/pkg/db"
	"github.com/Bio-OS/bioos/pkg/eventbus"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/server"
	"github.com/Bio-OS/bioos/pkg/storage"
)

const (
	ConfigFlagName = "config"
)

var cfgFile string

func init() {
	pflag.StringVarP(&cfgFile, ConfigFlagName, "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, YAML formats.")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("apiserver") // name of config file (without extension)
	if cfgFile != "" {               // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
		configDir := path.Dir(cfgFile)
		if configDir != "." {
			viper.AddConfigPath(configDir)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.AddConfigPath("conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

type Options struct {
	Client         *client.Options   `json:"client" mapstructure:"client"`
	ServerOption   *server.Options   `json:"server" mapstructure:"server"`
	LogOption      *log.Options      `json:"log" mapstructure:"log"`
	DBOption       *db.Options       `json:"db" mapstructure:"db"`
	AuthOption     *auth.Options     `json:"auth" mapstructure:"auth"`
	StorageOption  *storage.Options  `json:"storage" mapstructure:"storage"`
	EventBusOption *eventbus.Options `json:"eventBus" mapstructure:"eventBus"`
	WesOption      *wes.Options      `json:"wes" mapstructure:"wes"`
	NotebookOption *notebook.Options `json:"notebook" mapstructure:"notebook"`
}

func NewOptions() *Options {
	return &Options{
		ServerOption:   server.NewOptions(),
		LogOption:      log.NewOptions(),
		DBOption:       db.NewOptions(),
		AuthOption:     auth.NewOptions(),
		StorageOption:  storage.NewOptions(),
		EventBusOption: eventbus.NewOptions(),
		WesOption:      wes.NewOptions(),
		NotebookOption: notebook.NewOptions(),
	}
}

// Validate validate log options is valid.
func (o *Options) Validate() error {
	if err := o.ServerOption.Validate(); err != nil {
		return err
	}
	if err := o.LogOption.Validate(); err != nil {
		return err
	}
	if err := o.DBOption.Validate(); err != nil {
		return err
	}
	if err := o.AuthOption.Validate(); err != nil {
		return err
	}
	if err := o.EventBusOption.Validate(); err != nil {
		return err
	}
	if err := o.StorageOption.Validate(); err != nil {
		return err
	}
	if err := o.Client.Validate(); err != nil {
		return err
	}
	if err := o.WesOption.Validate(); err != nil {
		return err
	}
	if err := o.NotebookOption.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.ServerOption.AddFlags(fs)
	o.LogOption.AddFlags(fs)
	o.DBOption.AddFlags(fs)
	o.AuthOption.AddFlags(fs)
	o.StorageOption.AddFlags(fs)
	o.EventBusOption.AddFlags(fs)
	o.WesOption.AddFlags(fs)
	o.NotebookOption.AddFlags(fs)
}
