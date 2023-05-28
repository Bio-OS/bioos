package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type server struct {
	proto.UnimplementedNotebookServerServiceServer
	appService *application.Service
}

func NewServer(appService *application.Service) proto.NotebookServerServiceServer {
	return &server{
		appService: appService,
	}
}

func (s *server) RegisterServer(grpcServer grpc.ServiceRegistrar) {
	proto.RegisterNotebookServerServiceServer(grpcServer, s)
}

func (s *server) CreateNotebookServer(
	ctx context.Context, req *proto.CreateNotebookServerRequest) (*proto.CreateNotebookServerResponse, error) {
	id, err := s.appService.Commands.Create.Handle(ctx, newCreateCommand(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &proto.CreateNotebookServerResponse{
		Id: id,
	}, nil
}

func (s *server) GetNotebookServer(
	ctx context.Context, req *proto.GetNotebookServerRequest) (*proto.GetNotebookServerResponse, error) {
	nbsrv, err := s.appService.Queries.Get.Handle(ctx, newGetQuery(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return newGetNotebookServerResponse(nbsrv), nil
}

func (s *server) UpdateNotebookServerSettings(
	ctx context.Context, req *proto.UpdateNotebookServerSettingsRequest) (*proto.UpdateNotebookServerSettingsResponse, error) {
	err := s.appService.Commands.Update.Handle(ctx, newUpdateCommand(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &proto.UpdateNotebookServerSettingsResponse{}, nil
}

func (s *server) DeleteNotebookServer(
	ctx context.Context, req *proto.DeleteNotebookServerRequest) (*proto.DeleteNotebookServerResponse, error) {
	err := s.appService.Commands.Delete.Handle(ctx, newDeleteCommand(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &proto.DeleteNotebookServerResponse{}, nil
}

func (s *server) SwitchNotebookServer(
	ctx context.Context, req *proto.SwitchNotebookServerRequest) (*proto.SwitchNotebookServerResponse, error) {
	err := s.appService.Commands.Switch.Handle(ctx, newSwitchCommand(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &proto.SwitchNotebookServerResponse{}, nil
}

func (s *server) ListNotebookServers(
	ctx context.Context, req *proto.ListNotebookServersRequest) (*proto.ListNotebookServersResponse, error) {
	nbsrv, err := s.appService.Queries.List.Handle(ctx, newListQuery(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return newListNotebookServersResponse(nbsrv), nil
}
