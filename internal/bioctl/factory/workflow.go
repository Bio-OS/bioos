package factory

import (
	"context"

	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

type WorkflowClient interface {
	GetWorkflow(ctx context.Context, in *convert.GetWorkflowRequest) (*convert.GetWorkflowResponse, error)
	GetWorkflowVersion(ctx context.Context, in *convert.GetWorkflowVersionRequest) (*convert.GetWorkflowVersionResponse, error)
	ListWorkflowFiles(ctx context.Context, in *convert.ListWorkflowFilesRequest) (*convert.ListWorkflowFilesResponse, error)
	CreateWorkflow(ctx context.Context, in *convert.CreateWorkflowRequest) (*convert.CreateWorkflowResponse, error)
	DeleteWorkflow(ctx context.Context, in *convert.DeleteWorkflowRequest) (*convert.DeleteWorkflowResponse, error)
	UpdateWorkflow(ctx context.Context, in *convert.UpdateWorkflowRequest) (*convert.UpdateWorkflowResponse, error)
	ListWorkflow(ctx context.Context, in *convert.ListWorkflowsRequest) (*convert.ListWorkflowsResponse, error)
}

func (g *grpcClient) GetWorkflow(ctx context.Context, in *convert.GetWorkflowRequest) (*convert.GetWorkflowResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).GetWorkflow(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.GetWorkflowResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) GetWorkflowVersion(ctx context.Context, in *convert.GetWorkflowVersionRequest) (*convert.GetWorkflowVersionResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).GetWorkflowVersion(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.GetWorkflowVersionResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListWorkflowFiles(ctx context.Context, in *convert.ListWorkflowFilesRequest) (*convert.ListWorkflowFilesResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).ListWorkflowFiles(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListWorkflowFilesResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) CreateWorkflow(ctx context.Context, in *convert.CreateWorkflowRequest) (*convert.CreateWorkflowResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).CreateWorkflow(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.CreateWorkflowResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) DeleteWorkflow(ctx context.Context, in *convert.DeleteWorkflowRequest) (*convert.DeleteWorkflowResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).DeleteWorkflow(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteWorkflowResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) UpdateWorkflow(ctx context.Context, in *convert.UpdateWorkflowRequest) (*convert.UpdateWorkflowResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).UpdateWorkflow(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.UpdateWorkflowResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListWorkflow(ctx context.Context, in *convert.ListWorkflowsRequest) (*convert.ListWorkflowsResponse, error) {

	protoResp, err := workspaceproto.NewWorkflowServiceClient(g.conn).ListWorkflow(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListWorkflowsResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (h *httpClient) GetWorkflow(ctx context.Context, in *convert.GetWorkflowRequest) (*convert.GetWorkflowResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace-id}/workflow/{workflow-id}/file/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.GetWorkflowResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}

func (h *httpClient) GetWorkflowVersion(ctx context.Context, in *convert.GetWorkflowVersionRequest) (*convert.GetWorkflowVersionResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace-id}/workflow/{workflow-id}/version/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.GetWorkflowVersionResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}

func (h *httpClient) ListWorkflowFiles(ctx context.Context, in *convert.ListWorkflowFilesRequest) (*convert.ListWorkflowFilesResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace-id}/workflow/{workflow-id}/file"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListWorkflowFilesResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}

func (h *httpClient) CreateWorkflow(ctx context.Context, in *convert.CreateWorkflowRequest) (*convert.CreateWorkflowResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Post(h.url("workspace/{workspace-id}/workflow"))
	if err != nil {
		return nil, err
	}
	out := &convert.CreateWorkflowResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}

func (h *httpClient) DeleteWorkflow(ctx context.Context, in *convert.DeleteWorkflowRequest) (*convert.DeleteWorkflowResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Delete(h.url("workspace/{workspace-id}/workflow/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteWorkflowResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}

func (h *httpClient) UpdateWorkflow(ctx context.Context, in *convert.UpdateWorkflowRequest) (*convert.UpdateWorkflowResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Patch(h.url("workspace/{workspace-id}/workflow/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.UpdateWorkflowResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}

func (h *httpClient) ListWorkflow(ctx context.Context, in *convert.ListWorkflowsRequest) (*convert.ListWorkflowsResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace-id}/workflow"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListWorkflowsResponse{}
	err = convert.AssignFromHttpResponse(httpResp, out)
	return out, err
}
