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
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

type GrpcOption struct {
	Port string `json:"port" mapstructure:"port"`
	TLS  bool   `json:"tls" mapstructure:"tls"`
}

func NewGrpcOption() *GrpcOption {
	return &GrpcOption{}
}

func (o GrpcOption) Validate() error {
	return nil
}

func (o GrpcOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Port, "grpc-port", "50051", "grpc listen port")
	fs.BoolVar(&o.TLS, "grpc-tls-enable", false, "enable grpc tls")
}

type GRPCRegister func(grpc.ServiceRegistrar)

func GetGRPCRegister[T interface{}](register func(grpc.ServiceRegistrar, T), service T) GRPCRegister {
	return func(s grpc.ServiceRegistrar) {
		register(s, service)
	}
}
