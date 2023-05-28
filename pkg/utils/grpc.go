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
	"os"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
)

type ConnectInfo struct {
	// server address
	ServerAddr string `json:"serverAddr" mapstructure:"serverAddr"`
	// server name to override
	ServerName string `json:"serverName,omitempty" mapstructure:"serverName,omitempty"`
	// whether use tls
	Insecure bool `json:"insecure,omitempty" mapstructure:"insecure,omitempty"`
	// server cert file
	ServerCertFile string `json:"serverCertFile,omitempty" mapstructure:"serverCertFile,omitempty"`
	// client cert file
	ClientCertFile string `json:"clientCertFile,omitempty" mapstructure:"clientCertFile,omitempty"`
	// client  key file
	ClientCertKeyFile string `json:"clientCertKeyFile" mapstructure:"clientCertKeyFile,omitempty"`
	// ca file
	CaFile string `json:"caFile" mapstructure:"caFile,omitempty"`
}

func (c ConnectInfo) Validate() error {
	if c.ServerAddr == "" {
		return fmt.Errorf("server addr is empty")
	}
	return nil
}

// GrpcDial create grpc client connection
func GrpcDial(conn ConnectInfo, authInfo AuthInfo) (*grpc.ClientConn, error) {
	creds := insecure.NewCredentials()
	if conn.ServerCertFile != "" {
		var err error
		creds, err = credentials.NewClientTLSFromFile(conn.ServerCertFile, conn.ServerName)
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
			return nil, fmt.Errorf("failed to create TLS credentials %w", err)
		}
	} else if conn.ClientCertFile != "" && conn.ClientCertKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(conn.ClientCertFile, conn.ClientCertKeyFile)
		utils.CheckErr(err)
		certPool := x509.NewCertPool()
		if conn.CaFile != "" {
			ca, err := os.ReadFile(conn.CaFile)
			utils.CheckErr(err)
			certPool.AppendCertsFromPEM(ca)
		}

		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   conn.ServerName,
			RootCAs:      certPool,
		})
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(NewRPCCredentialFromAuthInfo(authInfo)),
	}

	return grpc.Dial(conn.ServerAddr, grpcOpts...)
}
