package wes

import (
	"fmt"

	"github.com/spf13/pflag"
)

const (
	DefaultEndpoint         = ""
	DefaultWDLBasePath      = "/api/ga4gh/wes/v1"
	DefaultNextflowBasePath = "/ga4gh/wes/v1"
	DefaultTimeout          = 10
	DefaultRetry            = 2
)

type BackendOptions struct {
	Endpoint string
	BasePath string
}

type Options struct {
	WDL      *BackendOptions
	Nextflow *BackendOptions
	Timeout  int
	Retry    int
}

// NewOptions new an event bus option.
func NewOptions() *Options {
	return &Options{
		WDL:      &BackendOptions{},
		Nextflow: &BackendOptions{},
	}
}

func (o Options) Validate() error {
	if o.Timeout < -1 {
		return fmt.Errorf("timeout second can not less than -1")
	}
	if o.Retry < 0 {
		return fmt.Errorf("retry times can not less than 0")
	}

	return nil
}

// AddFlags add event bus flags
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.WDL.Endpoint, "wes-wdl-endpoint", DefaultEndpoint, "wdl client endpoint")
	fs.StringVar(&o.WDL.BasePath, "wes-wdl-base-path", DefaultWDLBasePath, "wdl client api base path")
	fs.StringVar(&o.Nextflow.Endpoint, "wes-nextflow-endpoint", DefaultEndpoint, "nextflow client endpoint")
	fs.StringVar(&o.Nextflow.BasePath, "wes-nextflow-base-path", DefaultNextflowBasePath, "nextflow client api base path")
	fs.IntVar(&o.Timeout, "wes-timeout", DefaultTimeout, "wes client timeout")
	fs.IntVar(&o.Retry, "wes-retry", DefaultRetry, "wes client retry limit")
}
