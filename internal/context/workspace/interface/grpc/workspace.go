package grpc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/shaj13/go-guardian/v2/auth"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Bio-OS/bioos/internal/context/workspace/application"
	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workspace"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	pb "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// todo move datamodel seperately
type workspaceServer struct {
	pb.UnimplementedWorkspaceServiceServer
	pb.UnimplementedDataModelServiceServer
	workspaceService *application.WorkspaceService
}

// NewServer new a workspace rpc server.
func NewServer(workspaceService *application.WorkspaceService) pb.WorkspaceServiceServer {
	return &workspaceServer{
		workspaceService: workspaceService,
	}
}

// NewDataModelServer new a workspace rpc server of datamodel ...
func NewDataModelServer(workspaceService *application.WorkspaceService) pb.DataModelServiceServer {
	return &workspaceServer{
		workspaceService: workspaceService,
	}
}

func (s *workspaceServer) GetWorkspace(ctx context.Context, r *pb.GetWorkspaceRequest) (*pb.GetWorkspaceResponse, error) {
	log.Infow("GetWorkspace", "auth", auth.UserFromCtx(ctx))

	workspace, err := s.workspaceService.WorkspaceQueries.GetWorkspaceByID.Handle(ctx, &query.GetWorkspaceByIDQuery{
		ID: r.GetId(),
	})
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	storage := &pb.WorkspaceStorage{}
	if workspace.Storage.NFS != nil {
		storage.Nfs = &pb.NFSWorkspaceStorage{
			MountPath: workspace.Storage.NFS.MountPath,
		}
	}
	return &pb.GetWorkspaceResponse{
		Workspace: &pb.Workspace{
			Id:          workspace.ID,
			Name:        workspace.Name,
			Description: workspace.Description,
			CreatedAt:   timestamppb.New(workspace.CreatedAt),
			UpdatedAt:   timestamppb.New(workspace.UpdatedAt),
			Storage:     storage,
		},
	}, nil
}

func (s *workspaceServer) CreateWorkspace(ctx context.Context, r *pb.CreateWorkspaceRequest) (*pb.CreateWorkspaceResponse, error) {
	cmd := createWorkspaceVoToDto(r)
	id, err := s.workspaceService.WorkspaceCommands.CreateWorkspace.Handle(ctx, cmd)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.CreateWorkspaceResponse{
		Id: id,
	}, nil
}

func (s *workspaceServer) DeleteWorkspace(ctx context.Context, r *pb.DeleteWorkspaceRequest) (*pb.DeleteWorkspaceResponse, error) {
	cmd := &command.DeleteWorkspaceCommand{
		ID: r.GetId(),
	}
	err := s.workspaceService.WorkspaceCommands.DeleteWorkspace.Handle(ctx, cmd)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &pb.DeleteWorkspaceResponse{}, nil
}

func (s *workspaceServer) UpdateWorkspace(ctx context.Context, r *pb.UpdateWorkspaceRequest) (*pb.UpdateWorkspaceResponse, error) {
	var name *string
	if r.Name != "" {
		name = utils.PointString(r.GetName())
	}
	var description *string
	if r.Description != "" {
		description = utils.PointString(r.GetDescription())
	}
	cmd := &command.UpdateWorkspaceCommand{
		ID:          r.GetId(),
		Name:        name,
		Description: description,
	}
	err := s.workspaceService.WorkspaceCommands.UpdateWorkspace.Handle(ctx, cmd)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &pb.UpdateWorkspaceResponse{}, nil
}

func (s *workspaceServer) ListWorkspace(ctx context.Context, r *pb.ListWorkspaceRequest) (*pb.ListWorkspaceResponse, error) {
	listWorkspaceDto, err := listWorkspacesVoToDto(r)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	workspaceDtos, total, err := s.workspaceService.WorkspaceQueries.ListWorkspaces.Handle(ctx, listWorkspaceDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	if len(workspaceDtos) == 0 {
		return &pb.ListWorkspaceResponse{}, nil
	}

	workspaces := make([]*pb.Workspace, len(workspaceDtos))
	for i, ws := range workspaceDtos {
		workspaces[i] = workspaceItemDtoToVo(ws)
	}
	return &pb.ListWorkspaceResponse{
		Page:  r.Page,
		Size:  r.Size,
		Total: int32(total),
		Items: workspaces,
	}, nil
}

func (s *workspaceServer) ImportWorkspace(stream pb.WorkspaceService_ImportWorkspaceServer) error {
	var cmd *command.ImportWorkspaceCommand
	var r *pb.ImportWorkspaceRequest
	var err error
	var file *os.File
	var writer *bufio.Writer
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	for {
		r, err = stream.Recv()
		if cmd == nil {
			cmd = importWorkspaceVoToDto(r)
		}
		if err == io.EOF {
			err = s.workspaceService.WorkspaceCommands.ImportWorkspace.Handle(stream.Context(), cmd)
			if err != nil {
				return utils.ToGRPCError(err)
			}
			return stream.SendAndClose(&pb.ImportWorkspaceResponse{Id: cmd.ID})
		}
		if err != nil {
			return utils.ToGRPCError(err)
		}

		//In order to avoid the cost of wrong file writing, we validate file type here
		if path.Ext(r.GetFileName()) != consts.ImportWorkspaceFileTypeExt {
			return utils.ToGRPCError(fmt.Errorf("cannot parse [%s] in that we only support .zip file now", r.GetFileName()))
		}

		if file == nil {
			filePath := path.Join(cmd.Storage.NFS.MountPath, cmd.ID, cmd.FileName)
			err = os.MkdirAll(path.Dir(filePath), 0750)
			if err != nil {
				return err
			}

			file, err = os.Create(filepath.Clean(filePath))
			if err != nil {
				return err
			}

			writer = bufio.NewWriter(file)
		}
		_, err = writer.Write(r.GetContent())
		if err != nil {
			fmt.Println(err)
			if err != io.EOF {
				fmt.Println(err)
				return err
			}
		}
		writer.Flush()

		if err != nil {
			return utils.ToGRPCError(err)
		}
	}
}

func (s *workspaceServer) PatchDataModel(ctx context.Context, r *pb.PatchDataModelRequest) (*pb.PatchDataModelResponse, error) {
	patchDataModelDto := patchDataModelVoToDto(r)
	id, err := s.workspaceService.DataModelCommands.PatchDataModel.Handle(ctx, patchDataModelDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.PatchDataModelResponse{
		Id: id,
	}, nil
}

func (s *workspaceServer) DeleteDataModel(ctx context.Context, r *pb.DeleteDataModelRequest) (*pb.DeleteDataModelResponse, error) {
	deleteDataModelDto := deleteDataModelVoToDto(r)
	err := s.workspaceService.DataModelCommands.DeleteDataModel.Handle(ctx, deleteDataModelDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.DeleteDataModelResponse{}, nil
}

func (s *workspaceServer) ListDataModels(ctx context.Context, r *pb.ListDataModelsRequest) (*pb.ListDataModelsResponse, error) {
	listDataModelsDto := listDataModelsVoToDto(r)
	dataModels, err := s.workspaceService.DataModelQueries.ListDataModels.Handle(ctx, listDataModelsDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	items := make([]*pb.DataModel, 0, len(dataModels))
	for _, dataModel := range dataModels {
		items = append(items, dataModelsDtoToVo(dataModel))
	}

	return &pb.ListDataModelsResponse{
		Items: items,
	}, nil
}

func (s *workspaceServer) GetDataModel(ctx context.Context, r *pb.GetDataModelRequest) (*pb.GetDataModelResponse, error) {
	getDataModelDto := getDataModelVoToDto(r)
	dataModel, headers, err := s.workspaceService.DataModelQueries.GetDataModel.Handle(ctx, getDataModelDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.GetDataModelResponse{
		DataModel: dataModelsDtoToVo(dataModel),
		Headers:   headers,
	}, nil
}

func (s *workspaceServer) ListDataModelRows(ctx context.Context, r *pb.ListDataModelRowsRequest) (*pb.ListDataModelRowsResponse, error) {
	listDataModelRowsDto, err := listDataModelRowsVoToDto(r)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	headers, rows, total, err := s.workspaceService.DataModelQueries.ListDataModelRows.Handle(ctx, listDataModelRowsDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	rowVO := make([]*pb.Row, 0, len(rows))
	for _, row := range rows {
		rowVO = append(rowVO, rowDtoToVo(row))
	}

	return &pb.ListDataModelRowsResponse{
		Headers: headers,
		Rows:    rowVO,
		Page:    r.Page,
		Size:    r.Size,
		Total:   total,
	}, nil
}

func (s *workspaceServer) ListAllDataModelRowIDs(ctx context.Context, r *pb.ListAllDataModelRowIDsRequest) (*pb.ListAllDataModelRowIDsResponse, error) {
	listAllDataModelRowIDsDto := listAllDataModelRowIDsVoToDto(r)
	ids, err := s.workspaceService.DataModelQueries.ListAllDataModelRowIDs.Handle(ctx, listAllDataModelRowIDsDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &pb.ListAllDataModelRowIDsResponse{
		RowIDs: ids,
	}, nil
}
