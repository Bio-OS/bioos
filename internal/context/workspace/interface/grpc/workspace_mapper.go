package grpc

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	datamodelcommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/data-model"
	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workspace"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	pb "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func createWorkspaceVoToDto(req *pb.CreateWorkspaceRequest) *command.CreateWorkspaceCommand {
	return &command.CreateWorkspaceCommand{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Storage:     workspaceStorageVoToDto(req.Storage),
	}
}

func importWorkspaceVoToDto(req *pb.ImportWorkspaceRequest) *command.ImportWorkspaceCommand {
	return &command.ImportWorkspaceCommand{
		ID:       utils.GenWorkspaceID(),
		FileName: req.GetFileName(),
		Storage:  workspaceStorageVoToDto(req.Storage),
	}
}

func listWorkspacesVoToDto(req *pb.ListWorkspaceRequest) (*query.ListWorkspacesQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	return &query.ListWorkspacesQuery{
		Pg: *pg,
		Filter: &query.ListWorkspacesFilter{
			SearchWord: req.GetSearchWord(),
			Exact:      req.GetExact(),
			IDs:        req.Ids,
		},
	}, nil
}

func workspaceStorageVoToDto(s *pb.WorkspaceStorage) (res command.WorkspaceStorage) {
	if s.Nfs != nil {
		res.NFS = &command.NFSWorkspaceStorage{MountPath: s.Nfs.GetMountPath()}
	}
	return res
}

func workspaceItemDtoToVo(ws *query.WorkspaceItem) *pb.Workspace {
	return &pb.Workspace{
		Id:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description,
		Storage:     workspaceStorageDtoToVo(ws.Storage),
		CreatedAt:   timestamppb.New(ws.CreatedAt),
		UpdatedAt:   timestamppb.New(ws.UpdatedAt),
	}
}

func workspaceStorageDtoToVo(ws query.WorkspaceStorage) *pb.WorkspaceStorage {
	res := &pb.WorkspaceStorage{}
	if ws.NFS != nil {
		res.Nfs = &pb.NFSWorkspaceStorage{
			MountPath: ws.NFS.MountPath,
		}
	}
	return res
}

func patchDataModelVoToDto(req *pb.PatchDataModelRequest) *datamodelcommand.PatchDataModelCommand {
	rows := make([][]string, 0, len(req.Rows))
	for _, row := range req.Rows {
		rows = append(rows, rowVoToDto(row))
	}
	return &datamodelcommand.PatchDataModelCommand{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Async:       req.Async,
		Headers:     req.Headers,
		Rows:        rows,
	}
}

func deleteDataModelVoToDto(req *pb.DeleteDataModelRequest) *datamodelcommand.DeleteDataModelCommand {
	return &datamodelcommand.DeleteDataModelCommand{
		ID:          req.Id,
		WorkspaceID: req.WorkspaceID,
		Headers:     req.Headers,
		RowIDs:      req.RowIDs,
	}
}

func listDataModelsVoToDto(req *pb.ListDataModelsRequest) *datamodelquery.ListDataModelsQuery {
	return &datamodelquery.ListDataModelsQuery{
		WorkspaceID: req.WorkspaceID,
		Filter: &datamodelquery.ListDataModelsFilter{
			Types:      req.Types,
			SearchWord: req.SearchWord,
			Exact:      req.Exact,
			IDs:        req.Ids,
		},
	}
}

func getDataModelVoToDto(req *pb.GetDataModelRequest) *datamodelquery.GetDataModelQuery {
	return &datamodelquery.GetDataModelQuery{
		WorkspaceID: req.WorkspaceID,
		ID:          req.Id,
	}
}

func listDataModelRowsVoToDto(req *pb.ListDataModelRowsRequest) (*datamodelquery.ListDataModelRowsQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	if len(pg.Orders) > 1 {
		return nil, fmt.Errorf("only support one orderby")
	}
	return &datamodelquery.ListDataModelRowsQuery{
		WorkspaceID: req.WorkspaceID,
		ID:          req.Id,
		Pagination:  pg,
		Filter: &datamodelquery.ListDataModelRowsFilter{
			SearchWord: req.SearchWord,
			InSetIDs:   req.InSetIDs,
			RowIDs:     req.RowIDs,
		},
	}, nil
}

func listAllDataModelRowIDsVoToDto(req *pb.ListAllDataModelRowIDsRequest) *datamodelquery.ListAllDataModelRowIDsQuery {
	return &datamodelquery.ListAllDataModelRowIDsQuery{
		WorkspaceID: req.WorkspaceID,
		ID:          req.Id,
	}
}

func dataModelsDtoToVo(dataModel *datamodelquery.DataModel) *pb.DataModel {
	return &pb.DataModel{
		Id:       dataModel.ID,
		Name:     dataModel.Name,
		RowCount: dataModel.RowCount,
		Type:     dataModel.Type,
	}
}

func rowVoToDto(row *pb.Row) []string {
	res := make([]string, 0, len(row.Grids))
	for _, grid := range row.Grids {
		res = append(res, grid)
	}
	return res
}

func rowDtoToVo(row []string) *pb.Row {
	return &pb.Row{
		Grids: row,
	}
}
