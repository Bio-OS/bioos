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

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/client"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type DataModelClient interface {
	ListDataModels(context.Context, *workspaceproto.ListDataModelsRequest) (*workspaceproto.ListDataModelsResponse, error)
	GetDataModel(context.Context, *workspaceproto.GetDataModelRequest) (*workspaceproto.GetDataModelResponse, error)
	ListDataModelRows(context.Context, *workspaceproto.ListDataModelRowsRequest) (*workspaceproto.ListDataModelRowsResponse, error)
	ListAllDataModelRowIDs(context.Context, *workspaceproto.ListAllDataModelRowIDsRequest) (*workspaceproto.ListAllDataModelRowIDsResponse, error)
	PatchDataModel(context.Context, *workspaceproto.PatchDataModelRequest) (*workspaceproto.PatchDataModelResponse, error)
}

func NewDataModelClient(opts *client.Options) (DataModelClient, error) {
	if err := opts.Method.Validate(); err != nil {
		return nil, err
	}

	return dataModelClientImpl{
		opts: opts,
	}, nil
}

type dataModelClientImpl struct {
	opts *client.Options
}

func (d dataModelClientImpl) ListDataModels(ctx context.Context, req *workspaceproto.ListDataModelsRequest) (*workspaceproto.ListDataModelsResponse, error) {
	if d.opts.Method == client.GRPCMethod {
		conn, err := utils.GrpcDial(d.opts.ConnectInfo, d.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewDataModelServiceClient(conn)
		return client.ListDataModels(ctx, req)
	}
	return nil, fmt.Errorf("not support method")
}

func (d dataModelClientImpl) GetDataModel(ctx context.Context, req *workspaceproto.GetDataModelRequest) (*workspaceproto.GetDataModelResponse, error) {
	if d.opts.Method == client.GRPCMethod {
		conn, err := utils.GrpcDial(d.opts.ConnectInfo, d.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewDataModelServiceClient(conn)
		return client.GetDataModel(ctx, req)
	}
	return nil, fmt.Errorf("not support method")
}

func (d dataModelClientImpl) ListDataModelRows(ctx context.Context, req *workspaceproto.ListDataModelRowsRequest) (*workspaceproto.ListDataModelRowsResponse, error) {
	if d.opts.Method == client.GRPCMethod {
		conn, err := utils.GrpcDial(d.opts.ConnectInfo, d.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewDataModelServiceClient(conn)
		return client.ListDataModelRows(ctx, req)
	}
	return nil, fmt.Errorf("not support method")

}

func (d dataModelClientImpl) ListAllDataModelRowIDs(ctx context.Context, req *workspaceproto.ListAllDataModelRowIDsRequest) (*workspaceproto.ListAllDataModelRowIDsResponse, error) {
	if d.opts.Method == client.GRPCMethod {
		conn, err := utils.GrpcDial(d.opts.ConnectInfo, d.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewDataModelServiceClient(conn)
		return client.ListAllDataModelRowIDs(ctx, req)
	}
	return nil, fmt.Errorf("not support method")
}

func (d dataModelClientImpl) PatchDataModel(ctx context.Context, req *workspaceproto.PatchDataModelRequest) (*workspaceproto.PatchDataModelResponse, error) {
	if d.opts.Method == client.GRPCMethod {
		conn, err := utils.GrpcDial(d.opts.ConnectInfo, d.opts.AuthInfo)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := workspaceproto.NewDataModelServiceClient(conn)
		return client.PatchDataModel(ctx, req)
	}
	return nil, fmt.Errorf("not support method")
}
