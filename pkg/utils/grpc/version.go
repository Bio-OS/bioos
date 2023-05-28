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
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Bio-OS/bioos/internal/bioctl/utils"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/client"
	pkgutils "github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/version"
)

type VersionClient interface {
	Version(ctx context.Context) (*version.Info, error)
}

func NewVersionClient(opts *client.Options) (VersionClient, error) {
	if err := opts.Method.Validate(); err != nil {
		return nil, err
	}

	return &versionClientImpl{
		opts: opts,
	}, nil
}

type versionClientImpl struct {
	opts *client.Options
}

func (v versionClientImpl) Version(ctx context.Context) (*version.Info, error) {
	if v.opts.Method == client.HTTPMethod {
		cli, err := pkgutils.HttpDial(v.opts.ConnectInfo, v.opts.AuthInfo)
		if err != nil {
			utils.CheckErr(err)
		}
		var scheme = "https"
		if v.opts.Insecure {
			scheme = "http"
		}
		u := url.URL{
			Scheme: scheme,
			Host:   v.opts.ServerAddr,
			Path:   "version",
		}
		status, body, err := cli.Get(ctx, nil, u.String())
		if err != nil {
			return nil, err
		}
		if status != consts.StatusOK {
			return nil, fmt.Errorf("status code is %d", status)
		}
		var info *version.Info
		err = json.Unmarshal(body, &info)
		return info, err
	}
	conn, err := pkgutils.GrpcDial(v.opts.ConnectInfo, v.opts.AuthInfo)
	if err != nil {
		utils.CheckErr(err)
	}
	defer conn.Close()
	cli := workspaceproto.NewVersionServiceClient(conn)
	response, err := cli.Version(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return &version.Info{
		Version:      response.Version,
		GitBranch:    response.GitBranch,
		GitCommit:    response.GitCommit,
		GitTreeState: response.GitTreeState,
		BuildTime:    response.BuildTime,
		GoVersion:    response.GoVersion,
		Compiler:     response.Compiler,
		Platform:     response.Platform,
	}, nil
}
