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

package jupyterhub

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// API ...
type API struct {
	addr       string
	adminToken string
	rest       *resty.Client
}

// Server ...
type Server struct {
	Name             string    `json:"name"`
	Ready            bool      `json:"ready"`
	Pending          *string   `json:"pending"`
	Stopped          bool      `json:"stopped"`
	URL              string    `json:"url"`
	LastActivityTime time.Time `json:"last_activity"`
	ProgressURL      string    `json:"progress_url"`
	StartTime        time.Time `json:"started"`
}

// User ...
type User struct {
	Name    string
	Admin   bool
	Created time.Time
	Servers map[string]Server
}

// NewAPI ...
func NewAPI(hubAddress, adminToken string, hc *http.Client) *API {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &API{
		addr:       strings.TrimRight(hubAddress, "/"),
		adminToken: adminToken,
		rest:       resty.NewWithClient(hc),
	}
}

// Ping just test connection and authentication GetInfo
func (api *API) Ping(ctx context.Context) error {
	jerr := &Error{}
	resp, err := api.restR(ctx, jerr).Get(api.url("/info"))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return jerr
	}
	return nil
}

// CreateUser create a non-admin user into jupyterhub
func (api *API) CreateUser(ctx context.Context, username string) error {
	body := map[string]interface{}{
		"usernames": []string{username},
	}
	jerr := &Error{}
	resp, err := api.restR(ctx, jerr).SetBody(body).Post(api.url("/users"))
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusConflict {
		return ErrorConflict
	}
	if resp.IsError() {
		return jerr
	}
	return nil
}

// GetUser ...
func (api *API) GetUser(ctx context.Context, user string) (*User, error) {
	jerr := &Error{}
	info := &User{}
	resp, err := api.restR(ctx, jerr).SetResult(info).Get(api.url("/users/%s?include_stopped_servers", user))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrorNotFound
	}
	if resp.IsError() {
		return nil, jerr
	}
	return info, nil
}

// ListUsers ...
func (api *API) ListUsers(ctx context.Context) ([]User, error) {
	jerr := &Error{}
	var list []User
	resp, err := api.restR(ctx, jerr).SetResult(&list).Get(api.url("/users"))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, jerr
	}
	return list, nil
}

// StartServer ...
func (api *API) StartServer(ctx context.Context, user, servername string, options interface{}) error {
	jerr := &Error{}
	resp, err := api.restR(ctx, jerr).SetBody(options).Post(api.url("/users/%s/servers/%s", user, servername))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return jerr
	}
	return nil
}

// DeleteServer ...
func (api *API) DeleteServer(ctx context.Context, user, servername string) error {
	jerr := &Error{}
	resp, err := api.restR(ctx, jerr).SetBody(`{"remove":true}`).Delete(api.url("/users/%s/servers/%s", user, servername))
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return ErrorNotFound
	}
	if resp.IsError() {
		return jerr
	}
	return nil
}

// StopServer ...
func (api *API) StopServer(ctx context.Context, user, servername string) error {
	jerr := &Error{}
	resp, err := api.restR(ctx, jerr).Delete(api.url("/users/%s/servers/%s", user, servername))
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return ErrorNotFound
	}
	if resp.IsError() {
		return jerr
	}
	return nil
}

// DeleteUser ...
func (api *API) DeleteUser(ctx context.Context, user string) error {
	jerr := &Error{}
	resp, err := api.restR(ctx, jerr).Delete(api.url("/users/%s", user))
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return ErrorNotFound
	}
	if resp.IsError() {
		return jerr
	}
	return nil
}

func (api *API) restR(ctx context.Context, jerr *Error) *resty.Request {
	return api.rest.R().SetContext(ctx).SetError(jerr).SetAuthScheme("token").SetAuthToken(api.adminToken).ForceContentType("application/json")
}

func (api *API) url(apipath string, a ...interface{}) string {
	return url(api.addr, apipath, a...)
}

func url(hubAddress, apipath string, a ...interface{}) string {
	return strings.Join([]string{hubAddress, fmt.Sprintf(apipath, a...)}, "/hub/api")
}
