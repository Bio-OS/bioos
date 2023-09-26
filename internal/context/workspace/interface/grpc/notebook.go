package grpc

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type notebookServer struct {
	proto.UnimplementedNotebookServiceServer
	service *application.WorkspaceService
}

func NewNotebookServer(service *application.WorkspaceService) proto.NotebookServiceServer {
	return &notebookServer{
		service: service,
	}
}

func (s *notebookServer) CreateNotebook(ctx context.Context, req *proto.CreateNotebookRequest) (*proto.CreateNotebookResponse, error) {
	if err := s.service.NotebookCommands.Create.Handle(ctx, newNotebookCreateCommand(req)); err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &proto.CreateNotebookResponse{}, nil
}

func (s *notebookServer) DeleteNotebook(ctx context.Context, req *proto.DeleteNotebookRequest) (*proto.DeleteNotebookResponse, error) {
	if err := s.service.NotebookCommands.Delete.Handle(ctx, newNotebookDeleteCommand(req)); err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &proto.DeleteNotebookResponse{}, nil
}

func (s *notebookServer) ListNotebooks(ctx context.Context, req *proto.ListNotebooksRequest) (*proto.ListNotebooksResponse, error) {
	list, err := s.service.NotebookQueries.List.Handle(ctx, newNotebookListCommand(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return newListNotebooksResponse(list), nil
}

func (s *notebookServer) GetNotebook(ctx context.Context, req *proto.GetNotebookRequest) (*proto.GetNotebookResponse, error) {
	get, err := s.service.NotebookQueries.Get.Handle(ctx, newNotebookGetCommand(req))
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return newGetNotebookResponse(get), nil
}
