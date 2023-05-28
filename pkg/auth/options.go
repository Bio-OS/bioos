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

package auth

import (
	"github.com/spf13/pflag"

	"github.com/Bio-OS/bioos/pkg/auth/authn"
	"github.com/Bio-OS/bioos/pkg/auth/authz"
)

// Options stands for auth options.
type Options struct {
	AuthN *authn.Options `json:"authn,omitempty" mapstructure:"authn"`
	AuthZ *authz.Options `json:"authz,omitempty" mapstructure:"authz"`
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		AuthN: authn.NewOptions(),
		AuthZ: authz.NewOptions(),
	}
}

// Validate ...
func (o *Options) Validate() error {
	if err := o.AuthN.Validate(); err != nil {
		return err
	}
	if err := o.AuthZ.Validate(); err != nil {
		return err
	}
	return nil
}

// AddFlags ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.AuthN.AddFlags(fs)
	o.AuthZ.AddFlags(fs)
}
