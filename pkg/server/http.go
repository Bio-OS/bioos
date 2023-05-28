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
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/spf13/pflag"
)

type HttpOption struct {
	Port               string `json:"port" mapstructure:"port"`
	TLS                bool   `json:"tls" mapstructure:"tls"`
	MaxRequestBodySize int    `json:"maxRequestBodySize" mapstructure:"max-request-body-size"`
}

func NewHttpOption() *HttpOption {
	return &HttpOption{}
}

func (o HttpOption) Validate() error {
	return nil
}

func (o HttpOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Port, "http-port", "8888", "http listen port")
	fs.BoolVar(&o.TLS, "http-tls", false, "enable http tls")
	fs.IntVar(&o.MaxRequestBodySize, "http-max-request-body-size", 4*1024*1024, "http max request body size")
}

type RouteRegister interface {
	AddRoute(route.IRouter)
}
