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
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("JupyterHubAPI", func() {
	j := NewAPI("http://localhost:8081/jupyterhub", "abc", http.DefaultClient)

	ginkgo.AfterEach(func() {
		httpmock.Reset()
	})

	ginkgo.It("Ping", func() {
		httpmock.RegisterResponder(http.MethodGet, j.url("/info"), httpmock.NewStringResponder(http.StatusOK, ""))
		Expect(j.Ping(context.TODO())).ShouldNot(HaveOccurred())
	})

	ginkgo.It("CreateUser", func() {
		httpmock.RegisterResponder(http.MethodPost, j.url("/users"), httpmock.NewStringResponder(http.StatusOK, ""))
		Expect(j.CreateUser(context.TODO(), "haha")).ShouldNot(HaveOccurred())
	})

	ginkgo.It("GetUser", func() {
		httpmock.RegisterResponder(http.MethodGet, j.url(`/users/haha`), httpmock.NewStringResponder(200, `{"kind":"user","name":"12345","admin":false,"groups":[],"server":null,"pending":null,"created":"2022-02-25T08:13:47Z","last_activity":"2022-02-25T10:29:27.458641Z","servers":{"id":{"name":"id","last_activity":"2022-02-25T10:28:31.376000Z","started":"2022-02-25T10:21:37.457396Z","pending":null,"ready":true,"state":{"pod_name":"jupyter-12345--id"},"url":"/jupyterhub/user/12345/id/","user_options":{},"progress_url":"/jupyterhub/hub/api/users/12345/servers/id/progress"}},"auth_state":null}`))
		httpmock.RegisterResponder(http.MethodGet, j.url(`/users/yoyo`), httpmock.NewStringResponder(404, `{"status": 404, "message": "Not Found"}`))
		u, err := j.GetUser(context.TODO(), "haha")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(u.Name).To(Equal("12345"))
		Expect(u.Servers).To(HaveKey("id"))
		_, err = j.GetUser(context.TODO(), "yoyo")
		Expect(err).To(Equal(ErrorNotFound))
	})

	ginkgo.It("ListUsers", func() {
		httpmock.RegisterResponder(http.MethodGet, j.url("/users"), httpmock.NewStringResponder(200, `[{"kind": "user", "name": "van", "admin": true, "groups": [], "server": null, "pending": null, "created": "2022-02-09T09:36:57Z", "last_activity": "2022-03-21T14:04:01.787814Z", "servers": {}}, {"kind": "user", "name": "1000000000", "admin": false, "groups": [], "server": null, "pending": null, "created": "2022-02-26T03:29:12Z", "last_activity": "2022-03-20T12:32:01.073864Z", "servers": {}}, {"kind": "user", "name": "12345", "admin": false, "groups": [], "server": null, "pending": null, "created": "2022-03-11T15:05:09Z", "last_activity": "2022-03-11T16:31:06Z", "servers": {}}]`))
		u, err := j.ListUsers(context.TODO())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(u).To(HaveLen(3))
		Expect(u[0].Name).To(Equal("van"))
	})

	ginkgo.It("StartServer", func() {
		httpmock.RegisterResponder(http.MethodPost, "=~^"+j.url(`/users/(\S+)/servers/(\S+)`), httpmock.NewStringResponder(202, ""))
		Expect(j.StartServer(context.TODO(), "haha", "server", nil)).ShouldNot(HaveOccurred())
	})

	ginkgo.It("StartServer", func() {
		httpmock.RegisterResponder(http.MethodPost, "=~^"+j.url(`/users/(\S+)/servers/(\S+)`), httpmock.NewStringResponder(202, ""))
		Expect(j.StartServer(context.TODO(), "haha", "server", map[string]string{
			"profile": "test",
		})).ShouldNot(HaveOccurred())
	})

	ginkgo.It("DeleteServer", func() {
		httpmock.RegisterResponder(http.MethodDelete, "=~^"+j.url(`/users/(\S+)/servers/(\S+)`), httpmock.NewStringResponder(204, ""))
		Expect(j.DeleteServer(context.TODO(), "haha", "server")).ShouldNot(HaveOccurred())
	})

	ginkgo.It("StopServer", func() {
		httpmock.RegisterResponder(http.MethodDelete, "=~^"+j.url(`/users/(\S+)/servers/(\S+)`), httpmock.NewStringResponder(202, ""))
		Expect(j.StopServer(context.TODO(), "haha", "server")).ShouldNot(HaveOccurred())
	})

	ginkgo.It("DeleteUser", func() {
		httpmock.RegisterResponder(http.MethodDelete, "=~^"+j.url(`/users/(\S+)`), httpmock.NewStringResponder(204, ""))
		Expect(j.DeleteUser(context.TODO(), "haha")).ShouldNot(HaveOccurred())
	})
})
