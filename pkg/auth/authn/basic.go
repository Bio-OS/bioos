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
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
)

type BasicOption struct {
	TTL   int        `json:"ttl" mapstructure:"ttl"`
	Users BasicUsers `json:"users" mapstructure:"users"`
}

type BasicUsers []BasicUser

func (u *BasicUsers) Type() string {
	return "BasicUsers"
}

func (u *BasicUsers) String() string {
	return fmt.Sprintf("%v", []BasicUser(*u))
}

func (u *BasicUsers) Set(value string) error {
	var user BasicUser
	if err := json.Unmarshal([]byte(value), &user); err == nil {
		*u = append(*u, user)
	}
	return nil
}

type BasicUser struct {
	ID         string              `json:"ID" mapstructure:"ID"`
	Name       string              `json:"name" mapstructure:"name"`
	Password   string              `json:"password" mapstructure:"password"`
	Groups     []string            `json:"groups" mapstructure:"groups"`
	Extensions map[string][]string `json:"extensions" mapstructure:"extensions"`
}

func (u *BasicUser) Type() string {
	return "BasicUser"
}

func (u *BasicUser) String() string {
	return fmt.Sprintf("%v", *u)
}

func (u *BasicUser) Set(value string) error {
	var user BasicUser
	if err := json.Unmarshal([]byte(value), &user); err == nil {
		*u = user
	}
	return nil
}

func NewBasicOption() *BasicOption {
	return &BasicOption{}
}

func (o *BasicOption) Enabled() bool {
	return len(o.Users) > 0
}

func (o *BasicOption) Validate() error {
	return nil
}

func (o *BasicOption) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.TTL, "basic-ttl", 5, "cache ttl")
	fs.Var(&o.Users, "basic-users", "basic auth users")
}
