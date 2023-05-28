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

package grpc

import (
	"context"
	"net/http"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/shaj13/go-guardian/v2/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/middlewares"
)

func NewAuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(authFunc)
}

func NewAuthStreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(authFunc)
}

// GetCompatRequest ...
func GetCompatRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	// fix authorization
	if authorization, ok := md["authorization"]; ok {
		md["Authorization"] = authorization
	}
	applog.Infow("metadata", "md", md)
	req.Header = map[string][]string(md)
	return req, nil
}

func authFunc(ctx context.Context) (context.Context, error) {
	req, err := GetCompatRequest(ctx)
	if err != nil {
		return nil, err
	}
	user, err := middlewares.DefaultAuthenticator.Authenticate(ctx, req)
	applog.Debugw("AuthenticateRequest", "user", user, "err", err, "headers", req.Header, "method", req.Method, "path", req.RequestURI)
	if err == nil {
		return auth.CtxWithUser(ctx, user), nil
	}

	return nil, status.Errorf(codes.Unauthenticated, "fail to authenticate")
}
