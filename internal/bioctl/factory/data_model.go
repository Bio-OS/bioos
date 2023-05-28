package factory

import (
	"context"

	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

type DataModelClient interface {
	ListDataModels(ctx context.Context, in *convert.ListDataModelsRequest) (*convert.ListDataModelsResponse, error)
	GetDataModel(ctx context.Context, in *convert.GetDataModelRequest) (*convert.GetDataModelResponse, error)
	ListDataModelRows(ctx context.Context, in *convert.ListDataModelRowsRequest) (*convert.ListDataModelRowsResponse, error)
	PatchDataModel(ctx context.Context, in *convert.PatchDataModelRequest) (*convert.PatchDataModelResponse, error)
	DeleteDataModel(ctx context.Context, in *convert.DeleteDataModelRequest) (*convert.DeleteDataModelResponse, error)
	ListAllDataModelRowIDs(ctx context.Context, in *convert.ListAllDataModelRowIDsRequest) (*convert.ListAllDataModelRowIDsResponse, error)
}

func (g *grpcClient) ListDataModels(ctx context.Context, in *convert.ListDataModelsRequest) (*convert.ListDataModelsResponse, error) {

	protoResp, err := workspaceproto.NewDataModelServiceClient(g.conn).ListDataModels(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListDataModelsResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) GetDataModel(ctx context.Context, in *convert.GetDataModelRequest) (*convert.GetDataModelResponse, error) {

	protoResp, err := workspaceproto.NewDataModelServiceClient(g.conn).GetDataModel(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.GetDataModelResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListDataModelRows(ctx context.Context, in *convert.ListDataModelRowsRequest) (*convert.ListDataModelRowsResponse, error) {

	protoResp, err := workspaceproto.NewDataModelServiceClient(g.conn).ListDataModelRows(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListDataModelRowsResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) PatchDataModel(ctx context.Context, in *convert.PatchDataModelRequest) (*convert.PatchDataModelResponse, error) {

	protoResp, err := workspaceproto.NewDataModelServiceClient(g.conn).PatchDataModel(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.PatchDataModelResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) DeleteDataModel(ctx context.Context, in *convert.DeleteDataModelRequest) (*convert.DeleteDataModelResponse, error) {

	protoResp, err := workspaceproto.NewDataModelServiceClient(g.conn).DeleteDataModel(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteDataModelResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListAllDataModelRowIDs(ctx context.Context, in *convert.ListAllDataModelRowIDsRequest) (*convert.ListAllDataModelRowIDsResponse, error) {

	protoResp, err := workspaceproto.NewDataModelServiceClient(g.conn).ListAllDataModelRowIDs(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListAllDataModelRowIDsResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (h *httpClient) ListDataModels(ctx context.Context, in *convert.ListDataModelsRequest) (*convert.ListDataModelsResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/data_model"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListDataModelsResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) GetDataModel(ctx context.Context, in *convert.GetDataModelRequest) (*convert.GetDataModelResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/data_model/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.GetDataModelResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) ListDataModelRows(ctx context.Context, in *convert.ListDataModelRowsRequest) (*convert.ListDataModelRowsResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/data_model/{id}/rows"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListDataModelRowsResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) PatchDataModel(ctx context.Context, in *convert.PatchDataModelRequest) (*convert.PatchDataModelResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Patch(h.url("workspace/{workspace_id}/data_model"))
	if err != nil {
		return nil, err
	}
	out := &convert.PatchDataModelResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) DeleteDataModel(ctx context.Context, in *convert.DeleteDataModelRequest) (*convert.DeleteDataModelResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Delete(h.url("workspace/{workspace_id}/data_model/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteDataModelResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) ListAllDataModelRowIDs(ctx context.Context, in *convert.ListAllDataModelRowIDsRequest) (*convert.ListAllDataModelRowIDsResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace_id}/data_model/{id}/rows/ids"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListAllDataModelRowIDsResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}
