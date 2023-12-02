package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workflow"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	pb "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func createWorkflowVoToDto(req *pb.CreateWorkflowRequest) *command.CreateCommand {
	return &command.CreateCommand{
		Name:             req.GetName(),
		Description:      req.Description,
		WorkspaceID:      req.GetWorkspaceID(),
		Language:         req.GetLanguage(),
		Source:           req.GetSource(),
		Tag:              req.GetTag(),
		URL:              req.GetUrl(),
		Token:            req.GetToken(),
		MainWorkflowPath: req.GetMainWorkflowPath(),
		ID:               req.GetId(),
	}
}

func listWorkflowsVoToDto(req *pb.ListWorkflowRequest) (*query.ListQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	return &query.ListQuery{
		Pg: pg,
		Filter: &query.ListWorkflowsFilter{
			SearchWord: req.GetSearchWord(),
			IDs:        req.GetIds(),
			Exact:      req.GetExact(),
		},
	}, nil
}

func workflowItemDtoToVo(workflow *query.Workflow) *pb.Workflow {
	return &pb.Workflow{
		Id:            workflow.ID,
		Name:          workflow.Name,
		Description:   workflow.Description,
		LatestVersion: workflowVersionDtoToVo(workflow.LatestVersion),
		CreatedAt:     timestamppb.New(workflow.CreatedAt),
		UpdatedAt:     timestamppb.New(workflow.UpdatedAt),
	}
}

func workflowVersionDtoToVo(workflowVersion *query.WorkflowVersion) *pb.WorkflowVersion {
	inputs := make([]*pb.WorkflowParam, len(workflowVersion.Inputs))
	for index, input := range workflowVersion.Inputs {
		defVal := input.Default
		inputs[index] = &pb.WorkflowParam{
			Name:     input.Name,
			Type:     input.Type,
			Optional: input.Optional,
			Default:  &defVal,
		}
	}

	outputs := make([]*pb.WorkflowParam, len(workflowVersion.Outputs))
	for index, output := range workflowVersion.Outputs {
		defVal := output.Default
		outputs[index] = &pb.WorkflowParam{
			Name:     output.Name,
			Type:     output.Type,
			Optional: output.Optional,
			Default:  &defVal,
		}
	}
	files := make([]*pb.WorkflowFileInfo, len(workflowVersion.Files))
	for index, file := range workflowVersion.Files {
		files[index] = &pb.WorkflowFileInfo{
			Id:   file.ID,
			Path: file.Path,
		}
	}

	return &pb.WorkflowVersion{
		Id:               workflowVersion.ID,
		Status:           workflowVersion.Status,
		Message:          workflowVersion.Message,
		Language:         workflowVersion.Language,
		LanguageVersion:  workflowVersion.LanguageVersion,
		MainWorkflowPath: workflowVersion.MainWorkflowPath,
		Graph:            workflowVersion.Graph,
		Source:           workflowVersion.Source,
		Files:            files,
		Metadata:         workflowVersion.Metadata,
		Inputs:           inputs,
		Outputs:          outputs,
		CreatedAt:        timestamppb.New(workflowVersion.CreatedAt),
		UpdatedAt:        timestamppb.New(workflowVersion.UpdatedAt),
	}
}

func workflowFileDtoToVo(workflowFile *query.WorkflowFile) *pb.WorkflowFile {
	return &pb.WorkflowFile{
		Id:                workflowFile.ID,
		WorkflowVersionID: workflowFile.WorkflowVersionID,
		Path:              workflowFile.Path,
		Content:           workflowFile.Content,
		CreatedAt:         timestamppb.New(workflowFile.CreatedAt),
		UpdatedAt:         timestamppb.New(workflowFile.UpdatedAt),
	}
}

func listWorkflowFilesVoToDto(req *pb.ListWorkflowFilesRequest) (*query.ListFilesQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	ret := &query.ListFilesQuery{
		Pg: pg,
		Filter: &query.ListWorkflowFilesFilter{
			IDs: req.GetIds(),
		},
		WorkspaceID: req.GetWorkspaceID(),
		WorkflowID:  req.GetWorkflowID(),
	}
	if req.GetWorkflowVersionID() != "" {
		ret.WorkflowVersionID = req.GetWorkflowVersionID()
	}

	return ret, nil
}

func listWorkflowVersionsVoToDto(req *pb.ListWorkflowVersionsRequest) (*query.ListVersionsQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	return &query.ListVersionsQuery{
		Pg: pg,
		Filter: &query.ListWorkflowVersionsFilter{
			IDs: req.GetIds(),
		},
		WorkspaceID: req.GetWorkspaceID(),
		WorkflowID:  req.GetWorkflowID(),
	}, nil
}
