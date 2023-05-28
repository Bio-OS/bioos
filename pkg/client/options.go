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

package client

import (
	"fmt"

	pkgutils "github.com/Bio-OS/bioos/pkg/utils"
)

type Options struct {
	pkgutils.ConnectInfo `mapstructure:",squash"`
	pkgutils.AuthInfo    `mapstructure:",squash"`
	// whether use tls
	Insecure bool `json:"insecure,omitempty" mapstructure:"insecure,omitempty"`
	// connect method, grpc or http
	Method ConnectMethod `json:"method,omitempty" mapstructure:"method,omitempty"`
	// timeout seconds
	Timeout int `json:"timeout,omitempty" mapstructure:"timeout,omitempty"`
}

func (o Options) Validate() error {
	if err := o.Method.Validate(); err != nil {
		return err
	}
	if o.Timeout < -1 {
		return fmt.Errorf("timeout second can not less than -1")
	}
	if err := o.AuthInfo.Validate(); err != nil {
		return err
	}
	if err := o.ConnectInfo.Validate(); err != nil {
		return err
	}
	return nil
}

type ConnectMethod string

const (
	GRPCMethod ConnectMethod = "grpc"
	// HTTPMethod TODO  implement http method
	HTTPMethod ConnectMethod = "http"
)

func (m ConnectMethod) Validate() error {
	if m != GRPCMethod && m != HTTPMethod {
		return fmt.Errorf("connect method: %q invalid", m)
	}
	return nil
}
