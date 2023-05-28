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

type WorkflowClient interface {
	GetWorkflow(ctx context.Context, in *workspaceproto.GetWorkflowRequest) (*workspaceproto.GetWorkflowResponse, error)
	GetWorkflowVersion(ctx context.Context, in *workspaceproto.GetWorkflowVersionRequest) (*workspaceproto.GetWorkflowVersionResponse, error)
	ListWorkflowFiles(ctx context.Context, in *workspaceproto.ListWorkflowFilesRequest) (*workspaceproto.ListWorkflowFilesResponse, error)
	CreateWorkflow(ctx context.Context, in *workspaceproto.CreateWorkflowRequest) (*workspaceproto.CreateWorkflowResponse, error)
	DeleteWorkflow(ctx context.Context, in *workspaceproto.DeleteWorkflowRequest) (*workspaceproto.DeleteWorkflowResponse, error)
	UpdateWorkflow(ctx context.Context, in *workspaceproto.UpdateWorkflowRequest) (*workspaceproto.UpdateWorkflowResponse, error)
	ListWorkflow(ctx context.Context, in *workspaceproto.ListWorkflowRequest) (*workspaceproto.ListWorkflowResponse, error)
}

func NewWorkflowClient(opts *client.Options) (WorkflowClient, error) {
	if err := opts.Method.Validate(); err != nil {
		return nil, err
	}

	return workflowClientImpl{
		opts: opts,
	}, nil

}

var _ WorkflowClient = workflowClientImpl{}

type workflowClientImpl struct {
	opts *client.Options
}

func (w workflowClientImpl) GetWorkflow(ctx context.Context, in *workspaceproto.GetWorkflowRequest) (*workspaceproto.GetWorkflowResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.GetWorkflow(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workflowClientImpl) GetWorkflowVersion(ctx context.Context, in *workspaceproto.GetWorkflowVersionRequest) (*workspaceproto.GetWorkflowVersionResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.GetWorkflowVersion(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workflowClientImpl) ListWorkflowFiles(ctx context.Context, in *workspaceproto.ListWorkflowFilesRequest) (*workspaceproto.ListWorkflowFilesResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.ListWorkflowFiles(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workflowClientImpl) CreateWorkflow(ctx context.Context, in *workspaceproto.CreateWorkflowRequest) (*workspaceproto.CreateWorkflowResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.CreateWorkflow(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workflowClientImpl) DeleteWorkflow(ctx context.Context, in *workspaceproto.DeleteWorkflowRequest) (*workspaceproto.DeleteWorkflowResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.DeleteWorkflow(ctx, in)
	}
	return nil, fmt.Errorf("not support method")
}

func (w workflowClientImpl) UpdateWorkflow(ctx context.Context, in *workspaceproto.UpdateWorkflowRequest) (*workspaceproto.UpdateWorkflowResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.UpdateWorkflow(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}

func (w workflowClientImpl) ListWorkflow(ctx context.Context, in *workspaceproto.ListWorkflowRequest) (*workspaceproto.ListWorkflowResponse, error) {
	if w.opts.Method == client.GRPCMethod {
		conn, err := pkgutils.GrpcDial(w.opts.ConnectInfo, w.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		defer conn.Close()
		client := workspaceproto.NewWorkflowServiceClient(conn)
		return client.ListWorkflow(ctx, in)
	}

	return nil, fmt.Errorf("not support method")
}
