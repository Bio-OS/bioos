package wes

import (
	"fmt"

	"github.com/spf13/pflag"
)

const (
	DefaultEndpoint = ""
	DefaultBasePath = "/api/ga4gh/wes/v1"
	DefaultTimeout  = 10
	DefaultRetry    = 2
)

type Options struct {
	Endpoint string
	BasePath string
	Timeout  int
	Retry    int
}

// NewOptions new an event bus option.
func NewOptions() *Options {
	return &Options{}
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
	fs.StringVar(&o.Endpoint, "wes-endpoint", DefaultEndpoint, "wes client endpoint")
	fs.StringVar(&o.BasePath, "wes-base-path", DefaultBasePath, "wes client api base path")
	fs.IntVar(&o.Timeout, "wes-timeout", DefaultTimeout, "wes client timeout")
	fs.IntVar(&o.Retry, "wes-retry", DefaultRetry, "wes client retry limit")
}
