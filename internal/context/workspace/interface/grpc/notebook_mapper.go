package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/notebook"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

func newNotebookCreateCommand(req *proto.CreateNotebookRequest) *command.CreateCommand {
	return &command.CreateCommand{
		Name:        req.Name,
		WorkspaceID: req.WorkspaceID,
		Content:     req.Content,
	}
}

func newNotebookDeleteCommand(req *proto.DeleteNotebookRequest) *command.DeleteCommand {
	return &command.DeleteCommand{
		Name:        req.Name,
		WorkspaceID: req.WorkspaceID,
	}
}

func newNotebookListCommand(req *proto.ListNotebooksRequest) *query.ListQuery {
	return &query.ListQuery{
		WorkspaceID: req.WorkspaceID,
	}
}

func newNotebookVO(dto *query.Notebook) *proto.Notebook {
	return &proto.Notebook{
		Name:      dto.Name,
		Length:    dto.Size,
		UpdatedAt: timestamppb.New(dto.UpdateTime),
	}
}

func newListNotebooksResponse(list []*query.Notebook) *proto.ListNotebooksResponse {
	items := make([]*proto.Notebook, len(list))
	for i := range list {
		items[i] = newNotebookVO(list[i])
	}
	return &proto.ListNotebooksResponse{
		Items: items,
	}
}
