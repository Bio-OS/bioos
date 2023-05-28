package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/command"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

func convertResourceSizeToDTO(from *proto.ResourceSize, to *notebook.ResourceSize) {
	if from == nil || to == nil {
		return
	}
	if from.Gpu != nil {
		to.GPU = &notebook.GPU{
			Model:  from.Gpu.Model,
			Card:   from.Gpu.Card,
			Memory: from.Gpu.Memory,
		}
	}
	to.CPU = from.Cpu
	to.Memory = from.Memory
	to.Disk = from.Disk
}

func newResourceSizeVO(from *notebook.ResourceSize) *proto.ResourceSize {
	var gpu *proto.GPU
	if from.GPU != nil {
		gpu = &proto.GPU{
			Model:  from.GPU.Model,
			Card:   from.GPU.Card,
			Memory: from.GPU.Memory,
		}
	}
	return &proto.ResourceSize{
		Cpu:    from.CPU,
		Memory: from.Memory,
		Disk:   from.Disk,
		Gpu:    gpu,
	}
}

func newCreateCommand(req *proto.CreateNotebookServerRequest) *command.CreateCommand {
	res := &command.CreateCommand{
		WorkspaceID: req.WorkspaceID,
		Image:       req.Image,
	}
	convertResourceSizeToDTO(req.ResourceSize, &res.ResourceSize)
	return res
}

func newUpdateCommand(req *proto.UpdateNotebookServerSettingsRequest) *command.UpdateCommand {
	res := &command.UpdateCommand{
		ID: req.Id,
	}
	if len(req.Image) > 0 {
		res.Image = &req.Image
	}
	if req.ResourceSize != nil {
		res.ResourceSize = &notebook.ResourceSize{}
		convertResourceSizeToDTO(req.ResourceSize, res.ResourceSize)
	}
	return res
}

func newDeleteCommand(req *proto.DeleteNotebookServerRequest) *command.DeleteCommand {
	return &command.DeleteCommand{
		ID: req.Id,
	}
}

func newSwitchCommand(req *proto.SwitchNotebookServerRequest) *command.SwitchCommand {
	return &command.SwitchCommand{
		ID:    req.Id,
		OnOff: req.Onoff,
	}
}

func newGetQuery(req *proto.GetNotebookServerRequest) *query.GetQuery {
	return &query.GetQuery{
		ID:           req.Id,
		WorkspaceID:  req.WorkspaceID,
		EditNotebook: req.Notebook,
	}
}

func newGetNotebookServerResponse(dto *query.NotebookServer) *proto.GetNotebookServerResponse {
	return &proto.GetNotebookServerResponse{
		Id:           dto.ID,
		Image:        dto.Image,
		ResourceSize: newResourceSizeVO(&dto.ResourceSize),
		Status:       dto.Status,
		AccessURL:    dto.AccessURL,
		CreatedAt:    timestamppb.New(dto.CreateTime),
		UpdatedAt:    timestamppb.New(dto.UpdateTime),
	}
}

func newListQuery(req *proto.ListNotebookServersRequest) *query.ListQuery {
	return &query.ListQuery{
		WorkspaceID: req.WorkspaceID,
	}
}

func newNotebookServerVO(dto *query.NotebookServer) *proto.NotebookServer {
	return &proto.NotebookServer{
		Id:           dto.ID,
		Image:        dto.Image,
		ResourceSize: newResourceSizeVO(&dto.ResourceSize),
		Status:       dto.Status,
		CreatedAt:    timestamppb.New(dto.CreateTime),
		UpdatedAt:    timestamppb.New(dto.UpdateTime),
	}
}

func newListNotebookServersResponse(dto []query.NotebookServer) *proto.ListNotebookServersResponse {
	items := make([]*proto.NotebookServer, len(dto))
	for i := range dto {
		items[i] = newNotebookServerVO(&dto[i])
	}
	return &proto.ListNotebookServersResponse{
		Items: items,
	}
}
