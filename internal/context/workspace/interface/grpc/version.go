package grpc

import (
	"context"

	"github.com/shaj13/go-guardian/v2/auth"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/version"
)

type versionServer struct {
	pb.UnimplementedVersionServiceServer
}

// NewVersionServer new a version rpc server.
func NewVersionServer() pb.VersionServiceServer {
	return &versionServer{}
}

func (s *versionServer) RegisterServer(grpcServer grpc.ServiceRegistrar) {
	pb.RegisterVersionServiceServer(grpcServer, s)
}

func (s *versionServer) Version(ctx context.Context, _ *emptypb.Empty) (*pb.GetVersionResponse, error) {
	log.Infow("GetVersion", "auth", auth.UserFromCtx(ctx))

	info := version.Get()
	return &pb.GetVersionResponse{
		Version:      info.Version,
		GitBranch:    info.GitBranch,
		GitCommit:    info.GitCommit,
		GitTreeState: info.GitTreeState,
		BuildTime:    info.BuildTime,
		GoVersion:    info.GoVersion,
		Compiler:     info.Compiler,
		Platform:     info.Platform,
	}, nil
}
