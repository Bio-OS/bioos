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
	"fmt"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/client"
	pkgutils "github.com/Bio-OS/bioos/pkg/utils"
)

type WorkspaceClient interface {
	GetWorkspace(ctx context.Context, in *workspaceproto.GetWorkspaceRequest) (*workspaceproto.GetWorkspaceResponse, error)
	CreateWorkspace(ctx context.Context, in *workspaceproto.CreateWorkspaceRequest) (*workspaceproto.CreateWorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, in *workspaceproto.DeleteWorkspaceRequest) (*workspaceproto.DeleteWorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, in *workspaceproto.UpdateWorkspaceRequest) (*workspaceproto.UpdateWorkspaceResponse, error)
	ListWorkspace(ctx context.Context, in *workspaceproto.ListWorkspaceRequest) (*workspaceproto.ListWorkspaceResponse, error)
}

func NewWorkspaceClient(opts *client.Options) (WorkspaceClient, error) {
	if err := opts.Method.Validate(); err != nil {
		return nil, err
	}

	return workspaceClientImpl{
		opts: opts,
	}, nil

}

var _ WorkspaceClient = workspaceClientImpl{}

type workspaceClientImpl struct {
	opts *client.Options
}

func (w workspaceClientImpl) GetWorkspace(ctx context.Context, in *workspaceproto.GetWorkspaceRequest) (*workspaceproto.GetWorkspaceResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkspaceServiceClient(conn)
		return client.GetWorkspace(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workspaceClientImpl) CreateWorkspace(ctx context.Context, in *workspaceproto.CreateWorkspaceRequest) (*workspaceproto.CreateWorkspaceResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkspaceServiceClient(conn)
		return client.CreateWorkspace(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workspaceClientImpl) DeleteWorkspace(ctx context.Context, in *workspaceproto.DeleteWorkspaceRequest) (*workspaceproto.DeleteWorkspaceResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkspaceServiceClient(conn)
		return client.DeleteWorkspace(ctx, in)
	}
	return nil, fmt.Errorf("not support method")
}

func (w workspaceClientImpl) UpdateWorkspace(ctx context.Context, in *workspaceproto.UpdateWorkspaceRequest) (*workspaceproto.UpdateWorkspaceResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkspaceServiceClient(conn)
		return client.UpdateWorkspace(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workspaceClientImpl) ListWorkspace(ctx context.Context, in *workspaceproto.ListWorkspaceRequest) (*workspaceproto.ListWorkspaceResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkspaceServiceClient(conn)
		return client.ListWorkspace(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}
