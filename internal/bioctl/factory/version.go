package factory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"google.golang.org/protobuf/types/known/emptypb"

	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/version"
)

type VersionClient interface {
	Version(ctx context.Context) (*version.Info, error)
}

func (h *httpClient) Version(ctx context.Context) (*version.Info, error) {
	req := h.restR(ctx)
	httpResp, err := req.Get(h.url("version"))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() != consts.StatusOK {
		return nil, fmt.Errorf("status code is %d", httpResp.StatusCode())
	}
	var info *version.Info
	err = json.Unmarshal(httpResp.Body(), &info)
	return info, err

}

func (g *grpcClient) Version(ctx context.Context) (*version.Info, error) {
	cli := workspaceproto.NewVersionServiceClient(g.conn)
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
