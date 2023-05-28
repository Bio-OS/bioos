package grpc

import (
	"context"

	"github.com/shaj13/go-guardian/v2/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Bio-OS/bioos/internal/context/submission/application"
	runcommand "github.com/Bio-OS/bioos/internal/context/submission/application/command/run"
	command "github.com/Bio-OS/bioos/internal/context/submission/application/command/submission"
	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	pb "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc/proto"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type submissionServer struct {
	pb.UnimplementedSubmissionServiceServer
	submissionService *application.SubmissionService
}

// NewSubmissionServer new a submission rpc server.
func NewSubmissionServer(submissionService *application.SubmissionService) pb.SubmissionServiceServer {
	return &submissionServer{
		submissionService: submissionService,
	}
}

func (s *submissionServer) RegisterServer(grpcServer grpc.ServiceRegistrar) {
	pb.RegisterSubmissionServiceServer(grpcServer, s)
}

func (s *submissionServer) CheckSubmission(ctx context.Context, r *pb.CheckSubmissionRequest) (*pb.CheckSubmissionResponse, error) {
	applog.Infow("CheckSubmission", "auth", auth.UserFromCtx(ctx))

	isNameExist, err := s.submissionService.SubmissionQueries.Check.Handle(ctx, &query.CheckQuery{
		WorkspaceID: r.GetWorkspaceID(),
		Name:        r.GetName(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "get submission error:%v", err)
	}

	return &pb.CheckSubmissionResponse{
		IsNameExist: isNameExist,
	}, nil
}

func (s *submissionServer) ListSubmissions(ctx context.Context, r *pb.ListSubmissionsRequest) (*pb.ListSubmissionsResponse, error) {
	applog.Infow("ListSubmissions", "auth", auth.UserFromCtx(ctx))
	listSubmissionsDTO, err := listSubmissionsVOToDTO(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not convert list submissions req:%v", err)
	}
	items, count, err := s.submissionService.SubmissionQueries.List.Handle(ctx, listSubmissionsDTO)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "list submissions error:%v", err)
	}

	submissions := make([]*pb.SubmissionItem, len(items))
	for i, submission := range items {
		submissions[i] = submissionsItemDTOToVO(submission)
	}

	return &pb.ListSubmissionsResponse{
		Page:  r.Page,
		Size:  r.Size,
		Total: int32(count),
		Items: submissions,
	}, nil
}
func (s *submissionServer) CreateSubmission(ctx context.Context, r *pb.CreateSubmissionRequest) (*pb.CreateSubmissionResponse, error) {
	applog.Infow("CreateSubmission", "auth", auth.UserFromCtx(ctx))

	cmd := createSubmissionVOToDTO(r)
	id, err := s.submissionService.SubmissionCommands.CreateSubmission.Handle(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "create submissions error:%v", err)
	}
	return &pb.CreateSubmissionResponse{
		Id: id,
	}, nil
}
func (s *submissionServer) DeleteSubmission(ctx context.Context, r *pb.DeleteSubmissionRequest) (*pb.DeleteSubmissionResponse, error) {
	applog.Infow("DeleteSubmission", "auth", auth.UserFromCtx(ctx))

	err := s.submissionService.SubmissionCommands.DeleteSubmission.Handle(ctx, &command.DeleteSubmissionCommand{
		WorkspaceID: r.GetWorkspaceID(),
		ID:          r.GetId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "delete submissions error:%v", err)
	}
	return &pb.DeleteSubmissionResponse{}, nil
}
func (s *submissionServer) CancelSubmission(ctx context.Context, r *pb.CancelSubmissionRequest) (*pb.CancelSubmissionResponse, error) {
	applog.Infow("CancelSubmission", "auth", auth.UserFromCtx(ctx))

	err := s.submissionService.SubmissionCommands.CancelSubmission.Handle(ctx, &command.CancelSubmissionCommand{
		WorkspaceID: r.GetWorkspaceID(),
		ID:          r.GetId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "cancel submissions error:%v", err)
	}
	return &pb.CancelSubmissionResponse{}, nil
}
func (s *submissionServer) ListRuns(ctx context.Context, r *pb.ListRunsRequest) (*pb.ListRunsResponse, error) {
	applog.Infow("ListRuns", "auth", auth.UserFromCtx(ctx))

	listRunsDTO, err := listRunsVOToDTO(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not convert list runs req:%v", err)
	}

	items, count, err := s.submissionService.RunQueries.ListRuns.Handle(ctx, listRunsDTO)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "list runs error:%v", err)
	}

	submissions := make([]*pb.RunItem, len(items))
	for i, run := range items {
		submissions[i] = runItemDTOToVO(run)
	}

	return &pb.ListRunsResponse{
		Page:  r.Page,
		Size:  r.Size,
		Total: int32(count),
		Items: submissions,
	}, nil
}
func (s *submissionServer) CancelRun(ctx context.Context, r *pb.CancelRunRequest) (*pb.CancelRunResponse, error) {
	applog.Infow("CancelRun", "auth", auth.UserFromCtx(ctx))

	err := s.submissionService.RunCommands.CancelRun.Handle(ctx, &runcommand.CancelRunCommand{
		WorkspaceID:  r.GetWorkspaceID(),
		ID:           r.GetId(),
		SubmissionID: r.GetSubmissionID(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "cancel run error:%v", err)
	}
	return &pb.CancelRunResponse{}, nil
}
func (s *submissionServer) ListTasks(ctx context.Context, r *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	applog.Infow("ListTasks", "auth", auth.UserFromCtx(ctx))

	listTasksDTO, err := listTasksVOToDTO(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not convert list tasks req:%v", err)
	}

	items, count, err := s.submissionService.RunQueries.ListTasks.Handle(ctx, listTasksDTO)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "list tasks error:%v", err)
	}

	submissions := make([]*pb.TaskItem, len(items))
	for i, run := range items {
		submissions[i] = taskItemDTOToVO(run)
	}

	return &pb.ListTasksResponse{
		Page:  r.Page,
		Size:  r.Size,
		Total: int32(count),
		Items: submissions,
	}, nil
}
