//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/Bio-OS/bioos/pkg/utils"
)

// Options stands for server options.
type Options struct {
	Grpc        *GrpcOption `json:"grpc" mapstructure:"grpc"`
	Http        *HttpOption `json:"http" mapstructure:"http"`
	CertFile    string      `json:"certFile" mapstructure:"cert-file"`
	KeyFile     string      `json:"keyFile" mapstructure:"key-file"`
	CaFile      string      `json:"caFile" mapstructure:"ca-file"`
	WomtoolFile string      `json:"womtoolFile" mapstructure:"womtool-file"`
}

func (o Options) Validate() error {
	if err := o.Grpc.Validate(); err != nil {
		return err
	}
	if err := o.Http.Validate(); err != nil {
		return err
	}
	// validate womtool file
	if o.WomtoolFile == "" {
		return fmt.Errorf("womtool file can not be empty")
	}
	if err := utils.ValidateFileExist(o.WomtoolFile); err != nil {
		return err
	}
	return nil
}

func (o Options) AddFlags(fs *pflag.FlagSet) {
	o.Grpc.AddFlags(fs)
	o.Http.AddFlags(fs)
	fs.StringVar(&o.CertFile, "cert-file", "", "server cert file")
	fs.StringVar(&o.KeyFile, "key-file", "", "server key file")
	fs.StringVar(&o.CaFile, "ca-file", "", "ca file")
	fs.StringVar(&o.WomtoolFile, "womtool-file", "womtool.jar", "womtool file")
}

func NewOptions() *Options {
	return &Options{
		Grpc: NewGrpcOption(),
		Http: NewHttpOption(),
	}
}
