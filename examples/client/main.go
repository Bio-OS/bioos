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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/grpclog"

	pb "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

const (
	defaultName = ""
)

var (
	addr           = flag.String("addr", "localhost:50051", "the address to connect to")
	name           = flag.String("name", defaultName, "Name to query")
	certFile       = flag.String("cert-file", "", "cert file name")
	keyFile        = flag.String("key-file", "", "key file name")
	serverCertFile = flag.String("server-cert-file", "", "server cert file name")
	caFile         = flag.String("ca-file", "", "ca file name")
	serverName     = flag.String("server-name", "", "server name")
	username       = flag.String("username", "", "username")
	password       = flag.String("password", "", "password")
	accessToken    = flag.String("access-token", "", "oauth access token")
)

type basicAuth struct {
	username string
	password string
}

func (b basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	flag.Parse()
	creds := insecure.NewCredentials()
	if *serverName != "" {
		if *serverCertFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(*serverCertFile, *serverName)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else if *certFile != "" && *keyFile != "" {
			cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
			if err != nil {
				log.Fatalf("load cert err: %v", err)
			}
			certPool := x509.NewCertPool()
			if *caFile != "" {
				ca, _ := os.ReadFile(*caFile)
				certPool.AppendCertsFromPEM(ca)
			}

			creds = credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ServerName:   *serverName,
				RootCAs:      certPool,
			})
		}
	}

	var auth credentials.PerRPCCredentials

	grpcOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithTransportCredentials(creds),
	}
	if *username != "" && *password != "" {
		fmt.Printf("using basic auth: %s:%s\n", *username, *password)
		auth = &basicAuth{
			username: *username,
			password: *password,
		}
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(auth))
	} else if *accessToken != "" {
		auth = oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *accessToken})}
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(auth))
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpcOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewWorkspaceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetWorkspace(ctx, &pb.GetWorkspaceRequest{Id: *name})
	if err != nil {
		log.Fatalf("could not get workspace: %v", err)
	}
	log.Printf("Workspace: %s", r.GetWorkspace())
}
