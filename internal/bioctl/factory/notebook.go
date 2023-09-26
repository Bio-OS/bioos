package factory

import (
	"context"

	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

type NotebookClient interface {
	ListNotebooks(ctx context.Context, in *convert.ListNotebooksRequest) (*convert.ListNotebooksResponse, error)
	GetNotebook(ctx context.Context, in *convert.GetNotebookRequest) (*convert.GetNotebookResponse, error)
}

func (g *grpcClient) ListNotebooks(ctx context.Context, in *convert.ListNotebooksRequest) (*convert.ListNotebooksResponse, error) {

	protoResp, err := workspaceproto.NewNotebookServiceClient(g.conn).ListNotebooks(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListNotebooksResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) GetNotebook(ctx context.Context, in *convert.GetNotebookRequest) (*convert.GetNotebookResponse, error) {

	protoResp, err := workspaceproto.NewNotebookServiceClient(g.conn).GetNotebook(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.GetNotebookResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (h *httpClient) ListNotebooks(ctx context.Context, in *convert.ListNotebooksRequest) (*convert.ListNotebooksResponse, error) {

	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace-id}/notebook"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListNotebooksResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) GetNotebook(ctx context.Context, in *convert.GetNotebookRequest) (*convert.GetNotebookResponse, error) {

	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{workspace-id}/notebook/{name}"))
	if err != nil {
		return nil, err
	}
	out := &convert.GetNotebookResponse{}
	content, err := convert.RawBodyFromHttpResponse(httpResp)
	if err != nil {
		return nil, err
	}
	out.Content = content
	return out, nil
}
