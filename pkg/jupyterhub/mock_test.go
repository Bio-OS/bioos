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
	"testing"

	"github.com/onsi/gomega"
)

func TestMock(t *testing.T) {
	g := gomega.NewWithT(t)

	jupyterhubTestAddr := "http://hub.jupyterhub:8081/jupyterhub"
	m := NewJupyterHubMock(jupyterhubTestAddr, nil)

	c := m.Ping()
	g.Expect(c.url).To(gomega.Equal(jupyterhubTestAddr + "/hub/api/info"))

	c = m.GetUser(Any)
	g.Expect(c.url).To(gomega.Equal("=~^" + jupyterhubTestAddr + `/hub/api/users/(\S+)`))

	c = m.CreateUser(Any)
	g.Expect(c.url).To(gomega.Equal(jupyterhubTestAddr + "/hub/api/users"))

	c = m.ListUsers()
	g.Expect(c.url).To(gomega.Equal(jupyterhubTestAddr + "/hub/api/users"))

	c = m.CreateAPIToken(Any)
	g.Expect(c.url).To(gomega.Equal("=~^" + jupyterhubTestAddr + `/hub/api/users/(\S+)/tokens`))

	c = m.DeleteServer(Any, Any)
	g.Expect(c.url).To(gomega.Equal("=~^" + jupyterhubTestAddr + `/hub/api/users/(\S+)/servers/(\S+)`))

	c = m.StartServer(Any, Any)
	g.Expect(c.url).To(gomega.Equal("=~^" + jupyterhubTestAddr + `/hub/api/users/(\S+)/servers/(\S+)`))

	c = m.StopServer(Any, Any)
	g.Expect(c.url).To(gomega.Equal("=~^" + jupyterhubTestAddr + `/hub/api/users/(\S+)/servers/(\S+)`))

	c = m.DeleteUser(Any)
	g.Expect(c.url).To(gomega.Equal("=~^" + jupyterhubTestAddr + `/hub/api/users/(\S+)`))
}
