package factory

import (
	"context"

	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	submissionproto "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc/proto"
)

type SubmissionClient interface {
	ListSubmissions(ctx context.Context, in *convert.ListSubmissionsRequest) (*convert.ListSubmissionsResponse, error)
	CreateSubmission(ctx context.Context, in *convert.CreateSubmissionRequest) (*convert.CreateSubmissionResponse, error)
	DeleteSubmission(ctx context.Context, in *convert.DeleteSubmissionRequest) (*convert.DeleteSubmissionResponse, error)
	CancelSubmission(ctx context.Context, in *convert.CancelSubmissionRequest) (*convert.CancelSubmissionResponse, error)
	ListRuns(ctx context.Context, in *convert.ListRunsRequest) (*convert.ListRunsResponse, error)
	CancelRun(ctx context.Context, in *convert.CancelRunRequest) (*convert.CancelRunResponse, error)
	ListTasks(ctx context.Context, in *convert.ListTasksRequest) (*convert.ListTasksResponse, error)
}

func (g *grpcClient) ListSubmissions(ctx context.Context, in *convert.ListSubmissionsRequest) (*convert.ListSubmissionsResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).ListSubmissions(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListSubmissionsResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) CreateSubmission(ctx context.Context, in *convert.CreateSubmissionRequest) (*convert.CreateSubmissionResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).CreateSubmission(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.CreateSubmissionResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) DeleteSubmission(ctx context.Context, in *convert.DeleteSubmissionRequest) (*convert.DeleteSubmissionResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).DeleteSubmission(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteSubmissionResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) CancelSubmission(ctx context.Context, in *convert.CancelSubmissionRequest) (*convert.CancelSubmissionResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).CancelSubmission(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.CancelSubmissionResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListRuns(ctx context.Context, in *convert.ListRunsRequest) (*convert.ListRunsResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).ListRuns(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListRunsResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) CancelRun(ctx context.Context, in *convert.CancelRunRequest) (*convert.CancelRunResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).CancelRun(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.CancelRunResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListTasks(ctx context.Context, in *convert.ListTasksRequest) (*convert.ListTasksResponse, error) {

	protoResp, err := submissionproto.NewSubmissionServiceClient(g.conn).ListTasks(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListTasksResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (h *httpClient) ListSubmissions(ctx context.Context, in *convert.ListSubmissionsRequest) (*convert.ListSubmissionsResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/submission"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListSubmissionsResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) CreateSubmission(ctx context.Context, in *convert.CreateSubmissionRequest) (*convert.CreateSubmissionResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Post(h.url("workspace/{workspace_id}/submission"))
	if err != nil {
		return nil, err
	}
	out := &convert.CreateSubmissionResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) DeleteSubmission(ctx context.Context, in *convert.DeleteSubmissionRequest) (*convert.DeleteSubmissionResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Delete(h.url("workspace/{workspace_id}/submission/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteSubmissionResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) CancelSubmission(ctx context.Context, in *convert.CancelSubmissionRequest) (*convert.CancelSubmissionResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Post(h.url("workspace/{workspace_id}/submission/{id}/cancel"))
	if err != nil {
		return nil, err
	}
	out := &convert.CancelSubmissionResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) ListRuns(ctx context.Context, in *convert.ListRunsRequest) (*convert.ListRunsResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/submission/{submission_id}/run"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListRunsResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) CancelRun(ctx context.Context, in *convert.CancelRunRequest) (*convert.CancelRunResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Post(h.url("workspace/{workspace_id}/submission/{submission_id}/run/{id}/cancel"))
	if err != nil {
		return nil, err
	}
	out := &convert.CancelRunResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) ListTasks(ctx context.Context, in *convert.ListTasksRequest) (*convert.ListTasksResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/submission/{submission_id}/run/{run_id}/task"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListTasksResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}
