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

package authz

import (
	"github.com/spf13/pflag"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

type CasbinOption struct {
	ModelFile  string        `json:"model" mapstructure:"model"`
	Driver     string        `json:"driver,omitempty" mapstructure:"driver,omitempty"`
	MySQL      *MySQLOptions `json:"mysql,omitempty" mapstructure:"mysql,omitempty"`
	PolicyFile string        `json:"policy,omitempty" mapstructure:"policy,omitempty"`
}

func NewCabinOption() *CasbinOption {
	return &CasbinOption{
		MySQL: NewMySQLOptions(),
	}
}

func (o *CasbinOption) Validate() error {
	if o.ModelFile == "" {
		return apperrors.NewInvalidError("modelFile")
	}
	if err := o.MySQL.Validate(); err != nil {
		return apperrors.NewInvalidError()
	}
	return nil
}

func (o *CasbinOption) Enabled() bool {
	return o.ModelFile != "" && (o.PolicyFile != "" || o.Driver != "")
}

func (o *CasbinOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ModelFile, "casbin-model-file", "", "casbin model file")
	fs.StringVar(&o.PolicyFile, "casbin-policy-file", "", "casbin policy file")
	fs.StringVar(&o.Driver, "casbin-driver", "", "casbin driver: file or mysql")
	o.MySQL.AddFlags(fs)
}

type MySQLOptions struct {
	Username string `json:"username" mapstructure:"username,omitempty"`
	Password string `json:"password" mapstructure:"password,omitempty"`
	Host     string `json:"host" mapstructure:"host,omitempty"`
	Port     string `json:"port" mapstructure:"port,omitempty"`
}

// NewMySQLOptions new a mysql option.
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{}
}

// Validate validate log options is valid.
func (o *MySQLOptions) Validate() error {
	return nil
}

func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Username, "casbin-mysql-username", "", "casbin mysql db username")
	fs.StringVar(&o.Password, "casbin-mysql-password", "", "casbin mysql db password")
	fs.StringVar(&o.Host, "casbin-mysql-host", "localhost", "casbin mysql db host")
	fs.StringVar(&o.Port, "casbin-mysql-port", "3306", "casbin mysql db port")
}
