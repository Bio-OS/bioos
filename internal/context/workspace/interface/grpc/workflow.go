package grpc

import (
	"context"

	"github.com/shaj13/go-guardian/v2/auth"

	"github.com/Bio-OS/bioos/internal/context/workspace/application"
	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workflow"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	pb "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type workflowServer struct {
	pb.UnimplementedWorkflowServiceServer
	workflowService *application.WorkspaceService
}

// NewWorkflowServer new a workflow rpc server.
func NewWorkflowServer(workflowService *application.WorkspaceService) pb.WorkflowServiceServer {
	return &workflowServer{
		workflowService: workflowService,
	}
}

func (s *workflowServer) GetWorkflow(ctx context.Context, r *pb.GetWorkflowRequest) (*pb.GetWorkflowResponse, error) {
	log.Infow("GetWorkflow", "auth", auth.UserFromCtx(ctx))

	workflow, err := s.workflowService.WorkflowQueries.GetByID.Handle(ctx, &query.GetQuery{
		ID:          r.GetId(),
		WorkspaceID: r.GetWorkspaceID(),
	})
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.GetWorkflowResponse{
		Workflow: workflowItemDtoToVo(workflow),
	}, nil
}

func (s *workflowServer) CreateWorkflow(ctx context.Context, r *pb.CreateWorkflowRequest) (*pb.CreateWorkflowResponse, error) {
	cmd := createWorkflowVoToDto(r)
	id, err := s.workflowService.WorkflowCommands.Create.Handle(ctx, cmd)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.CreateWorkflowResponse{
		Id: id,
	}, nil
}

func (s *workflowServer) DeleteWorkflow(ctx context.Context, r *pb.DeleteWorkflowRequest) (*pb.DeleteWorkflowResponse, error) {
	cmd := &command.DeleteCommand{
		ID:          r.GetId(),
		WorkspaceID: r.GetWorkspaceID(),
	}
	err := s.workflowService.WorkflowCommands.Delete.Handle(ctx, cmd)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &pb.DeleteWorkflowResponse{}, nil
}

func (s *workflowServer) UpdateWorkflow(ctx context.Context, r *pb.UpdateWorkflowRequest) (*pb.UpdateWorkflowResponse, error) {
	cmd := &command.UpdateCommand{
		ID:               r.Id,
		WorkspaceID:      r.WorkspaceID,
		Name:             r.Name,
		Description:      r.Description,
		Language:         r.Language,
		Source:           r.Source,
		URL:              r.Url,
		Tag:              r.Tag,
		Token:            r.Token,
		MainWorkflowPath: r.MainWorkflowPath,
	}
	err := s.workflowService.WorkflowCommands.Update.Handle(ctx, cmd)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	return &pb.UpdateWorkflowResponse{}, nil
}

func (s *workflowServer) ListWorkflow(ctx context.Context, r *pb.ListWorkflowRequest) (*pb.ListWorkflowResponse, error) {
	listWorkflowDto, err := listWorkflowsVoToDto(r)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	workflowDtos, total, err := s.workflowService.WorkflowQueries.ListWorkflows.Handle(ctx, listWorkflowDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	if len(workflowDtos) == 0 {
		return &pb.ListWorkflowResponse{}, nil
	}

	workflows := make([]*pb.Workflow, len(workflowDtos))
	for i, ws := range workflowDtos {
		workflows[i] = workflowItemDtoToVo(ws)
	}
	return &pb.ListWorkflowResponse{
		Page:  r.Page,
		Size:  r.Size,
		Total: int32(total),
		Items: workflows,
	}, nil
}

func (s *workflowServer) ListWorkflowFiles(ctx context.Context, r *pb.ListWorkflowFilesRequest) (*pb.ListWorkflowFilesResponse, error) {
	log.Infow("ListWorkflowFiles", "auth", auth.UserFromCtx(ctx))
	listWorkflowFilesDto, err := listWorkflowFilesVoToDto(r)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	workflowFilesDtos, total, err := s.workflowService.WorkflowQueries.ListFiles.Handle(ctx, listWorkflowFilesDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	if len(workflowFilesDtos) == 0 {
		return &pb.ListWorkflowFilesResponse{}, nil
	}

	workflowFiles := make([]*pb.WorkflowFile, len(workflowFilesDtos))
	for i, ws := range workflowFilesDtos {
		workflowFiles[i] = workflowFileDtoToVo(ws)
	}
	return &pb.ListWorkflowFilesResponse{
		Page:        r.Page,
		Size:        r.Size,
		Total:       int32(total),
		WorkspaceID: r.GetWorkspaceID(),
		WorkflowID:  r.GetWorkflowID(),
		Files:       workflowFiles,
	}, nil
}

func (s *workflowServer) GetWorkflowVersion(ctx context.Context, r *pb.GetWorkflowVersionRequest) (*pb.GetWorkflowVersionResponse, error) {
	log.Infow("GetWorkflowVersion", "auth", auth.UserFromCtx(ctx))

	workflowVersion, err := s.workflowService.WorkflowQueries.GetVersion.Handle(ctx, &query.GetVersionQuery{
		ID:          r.GetId(),
		WorkflowID:  r.GetWorkflowID(),
		WorkspaceID: r.GetWorkspaceID(),
	})
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.GetWorkflowVersionResponse{
		Version: workflowVersionDtoToVo(workflowVersion),
	}, nil
}

func (s *workflowServer) ListWorkflowVersions(ctx context.Context, r *pb.ListWorkflowVersionsRequest) (*pb.ListWorkflowVersionsResponse, error) {
	log.Infow("ListWorkflowVersions", "auth", auth.UserFromCtx(ctx))
	listWorkflowVersionsDto, err := listWorkflowVersionsVoToDto(r)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	workflowVersionsDtos, total, err := s.workflowService.WorkflowQueries.ListVersions.Handle(ctx, listWorkflowVersionsDto)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}
	if len(workflowVersionsDtos) == 0 {
		return &pb.ListWorkflowVersionsResponse{}, nil
	}

	workflowVersions := make([]*pb.WorkflowVersion, len(workflowVersionsDtos))
	for i, wv := range workflowVersionsDtos {
		workflowVersions[i] = workflowVersionDtoToVo(wv)
	}
	return &pb.ListWorkflowVersionsResponse{
		Page:        r.Page,
		Size:        r.Size,
		Total:       int32(total),
		Items:       workflowVersions,
		WorkspaceID: r.GetWorkspaceID(),
		WorkflowID:  r.GetWorkflowID(),
	}, nil
}

func (s *workflowServer) GetWorkflowFile(ctx context.Context, r *pb.GetWorkflowFileRequest) (*pb.GetWorkflowFileResponse, error) {
	log.Infow("GetWorkflowFile", "auth", auth.UserFromCtx(ctx))

	workflowFile, err := s.workflowService.WorkflowQueries.GetFile.Handle(ctx, &query.GetFileQuery{
		ID:          r.GetId(),
		WorkflowID:  r.GetWorkflowID(),
		WorkspaceID: r.GetWorkspaceID(),
	})
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.GetWorkflowFileResponse{
		File: workflowFileDtoToVo(workflowFile),
	}, nil
}
