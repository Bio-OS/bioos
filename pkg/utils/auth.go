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

package utils

import (
	"context"
	"encoding/base64"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

type AuthInfo struct {
	// username
	Username string `json:"username,omitempty" mapstructure:"username,omitempty"`
	// password
	Password string `json:"password,omitempty" mapstructure:"password,omitempty"`
	// authorization token
	AuthToken string `json:"authToken,omitempty" mapstructure:"authToken,omitempty"`
}

func (a AuthInfo) Validate() error {
	if a.Username == "" && a.Password == "" && a.AuthToken == "" {
		return fmt.Errorf("basic auth and token can not be empty")
	}
	return nil

}
func NewRPCCredentialFromAuthInfo(authInfo AuthInfo) credentials.PerRPCCredentials {
	if authInfo.Username != "" {
		return NewBasicAuthRPCCredentials(authInfo.Username, authInfo.Password)
	}
	return NewTokenRPCCredentials(authInfo.AuthToken)
}

// NewTokenRPCCredentials create rpc credentials with token
func NewTokenRPCCredentials(token string) credentials.PerRPCCredentials {
	return oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
}

// NewBasicAuthRPCCredentials create rpc credentials with basic auth
func NewBasicAuthRPCCredentials(username, password string) credentials.PerRPCCredentials {
	return basicAuth{
		username: username,
		password: password,
	}
}

type basicAuth struct {
	username string
	password string
}

func (b basicAuth) getEncodeCode() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		b.username, b.password)))
}

func (b basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + b.getEncodeCode(),
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}
