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

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

type JWTOption struct {
	ID        string `json:"id" mapstructure:"id"`
	Secret    string `json:"secret" mapstructure:"secret"`
	Algorithm string `json:"algorithm" mapstructure:"algorithm"`
}

func NewJWTOption() *JWTOption {
	return &JWTOption{
		ID:        "secret-id",
		Secret:    "secret",
		Algorithm: "HS256",
	}
}

func (o *JWTOption) Enabled() bool {
	return o.Secret != "" && o.Algorithm != ""
}

func (o *JWTOption) Validate() error {
	if o.Algorithm == "" {
		return apperrors.NewInvalidError("algorithm")
	}

	return nil
}

func (o *JWTOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ID, "jwt-id", "", "jwt id")
	fs.StringVar(&o.Secret, "jwt-secret", "", "jwt secret")
	fs.StringVar(&o.Algorithm, "jwt-algorithm", "", "jwt algorithm")
}
