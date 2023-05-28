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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	c "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
)

// HttpDial create grpc client connection
func HttpDial(conn ConnectInfo, authInfo AuthInfo) (*c.Client, error) {
	var options []config.ClientOption
	options = append(options, c.WithDialTimeout(consts.DefaultDialTimeout))

	if !conn.Insecure {
		if conn.ClientCertFile != "" && conn.ClientCertKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(conn.ClientCertFile, conn.ClientCertKeyFile)
			utils.CheckErr(err)
			certPool := x509.NewCertPool()
			if conn.CaFile != "" {
				ca, err := os.ReadFile(conn.CaFile)
				utils.CheckErr(err)
				certPool.AppendCertsFromPEM(ca)
			}
			cfg := &tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientCAs:    certPool,
			}
			options = append(options, c.WithTLSConfig(cfg))
		}
	}
	return c.NewClient(
		options...,
	)
}

type TransportWithAuth struct {
	basicAuth
	*http.Transport
}

func NewTransportWithAuth(conn ConnectInfo, authInfo AuthInfo) (*TransportWithAuth, error) {
	transport := http.DefaultTransport.(*http.Transport)
	if !conn.Insecure {
		conf, err := NewHTTPTlSConfig(conn)
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig = conf
	}
	return &TransportWithAuth{
		basicAuth{
			username: authInfo.Username,
			password: authInfo.Password,
		},
		transport,
	}, nil
}

func (t TransportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", t.basicAuth.getEncodeCode()))
	return t.Transport.RoundTrip(req)
}

func (t TransportWithAuth) Client() *http.Client {
	return &http.Client{Transport: t}
}

func NewHTTPTlSConfig(conn ConnectInfo) (*tls.Config, error) {
	var cfg *tls.Config
	if conn.ClientCertFile != "" && conn.ClientCertKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(conn.ClientCertFile, conn.ClientCertKeyFile)
		utils.CheckErr(err)
		certPool := x509.NewCertPool()
		if conn.CaFile != "" {
			ca, err := os.ReadFile(conn.CaFile)
			utils.CheckErr(err)
			certPool.AppendCertsFromPEM(ca)
		}
		cfg = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   conn.ServerName,
			RootCAs:      certPool,
		}
	}
	return cfg, nil
}
