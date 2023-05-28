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
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

const (
	// Any ...
	Any = `(\S+)`
)

// Call ...
type Call struct {
	method string
	url    string
}

// Success ...
func (c *Call) Success(a ...interface{}) {
	c.Return(200, a...)
}

// Return ...
func (c *Call) Return(code int, a ...interface{}) {
	var resp httpmock.Responder
	if len(a) > 0 {
		resp, _ = httpmock.NewJsonResponder(code, a[0])
	} else {
		resp = httpmock.NewStringResponder(code, "")
	}
	httpmock.RegisterResponder(c.method, c.url, resp)
}

// Mock ...
type Mock struct {
	addr string
}

// NewJupyterHubMock hubAddress must be same witch API.addr
func NewJupyterHubMock(hubAddress string, hc *http.Client) *Mock {
	if hc == nil {
		httpmock.ActivateNonDefault(http.DefaultClient)
	} else {
		httpmock.ActivateNonDefault(hc)
	}
	return &Mock{
		addr: strings.TrimRight(hubAddress, "/"),
	}
}

// Close call it when test finish
func (m *Mock) Close() {
	httpmock.DeactivateAndReset()
}

// Ping ...
func (m *Mock) Ping() *Call {
	return &Call{
		method: http.MethodGet,
		url:    m.url("/info"),
	}
}

// CreateUser ...
func (m *Mock) CreateUser(username string) *Call {
	return &Call{
		method: http.MethodPost,
		url:    m.url("/users"),
	}
}

// GetUser ...
func (m *Mock) GetUser(username string) *Call {
	return &Call{
		method: http.MethodGet,
		url:    m.url("/users/%s", username),
	}
}

// ListUsers ...
func (m *Mock) ListUsers() *Call {
	return &Call{
		method: http.MethodGet,
		url:    m.url("/users"),
	}
}

// CreateAPIToken ...
func (m *Mock) CreateAPIToken(username string) *Call {
	return &Call{
		method: http.MethodPost,
		url:    m.url("/users/%s/tokens", username),
	}
}

// StartServer ...
func (m *Mock) StartServer(username, servername string, options ...interface{}) *Call {
	return &Call{
		method: http.MethodPost,
		url:    m.url("/users/%s/servers/%s", username, servername),
	}
}

// DeleteServer ...
func (m *Mock) DeleteServer(username, servername string) *Call {
	return &Call{
		method: http.MethodDelete,
		url:    m.url("/users/%s/servers/%s", username, servername),
	}
}

// StopServer ...
func (m *Mock) StopServer(username, servername string) *Call {
	return &Call{
		method: http.MethodDelete,
		url:    m.url("/users/%s/servers/%s", username, servername),
	}
}

// DeleteUser ...
func (m *Mock) DeleteUser(username string) *Call {
	return &Call{
		method: http.MethodDelete,
		url:    m.url("/users/%s", username),
	}
}

// Reset ...
func (m *Mock) Reset() {
	httpmock.Reset()
}

func (m *Mock) url(apipath string, a ...interface{}) string {
	u := url(m.addr, apipath, a...)
	for _, s := range a {
		if s == Any {
			u = "=~^" + u
			break
		}
	}
	return u
}
