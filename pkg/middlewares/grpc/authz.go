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
	"strings"

	"github.com/shaj13/go-guardian/v2/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/middlewares"
)

// RBACUnaryServerChain check rbac permission in unary.
func RBACUnaryServerChain() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err := checkPermission(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// RBACStreamServerChain check rbac permission in stream.
func RBACStreamServerChain() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := checkPermission(ss.Context(), info.FullMethod); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func checkPermission(ctx context.Context, fullMethod string) error {
	FullMethod := strings.Split(fullMethod, "/")
	if len(FullMethod) < 2 {
		return grpc.Errorf(codes.Internal, "not enough params in full method")
	}
	obj := FullMethod[1]    // proto.WorkspaceService
	action := FullMethod[2] // GetWorkspace
	user := auth.UserFromCtx(ctx)
	allowed, err := middlewares.DefaultAuthorizer.Authorize(user.GetUserName(), obj, action)
	if err != nil {
		return grpc.Errorf(codes.Internal, "Authorizer internal error")
	}
	if !allowed {
		applog.Debugw("user not allowed", "obj", obj, "action", action, "user", user)
		return grpc.Errorf(codes.Unauthenticated, "Permission denied")
	}

	applog.Debugw("user permitted to enter", "obj", obj, "action", action, "user", user)
	return nil
}
