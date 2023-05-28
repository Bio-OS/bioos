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

package authn

import (
	"github.com/spf13/pflag"
)

// Options stands for auth options.
type Options struct {
	Basic *BasicOption `json:"basic,omitempty" mapstructure:"basic"`
	JWT   *JWTOption   `json:"jwt,omitempty" mapstructure:"jwt"`
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		Basic: NewBasicOption(),
		JWT:   NewJWTOption(),
	}
}

// Validate ...
func (o *Options) Validate() error {
	if o.Basic.Enabled() {
		if err := o.Basic.Validate(); err != nil {
			return err
		}
	}
	if o.JWT.Enabled() {
		if err := o.JWT.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// AddFlags ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Basic.AddFlags(fs)
	o.JWT.AddFlags(fs)
}
