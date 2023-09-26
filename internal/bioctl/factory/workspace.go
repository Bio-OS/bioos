package factory

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/Bio-OS/bioos/internal/bioctl/factory/convert"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
)

type WorkspaceClient interface {
	CreateWorkspace(ctx context.Context, in *convert.CreateWorkspaceRequest) (*convert.CreateWorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, in *convert.DeleteWorkspaceRequest) (*convert.DeleteWorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, in *convert.UpdateWorkspaceRequest) (*convert.UpdateWorkspaceResponse, error)
	ListWorkspaces(ctx context.Context, in *convert.ListWorkspacesRequest) (*convert.ListWorkspacesResponse, error)
	ImportWorkspace(ctx context.Context, in *convert.ImportWorkspaceRequest) (*convert.ImportWorkspaceResponse, error)
	GetWorkspace(ctx context.Context, in *convert.GetWorkspaceRequest) (*convert.GetWorkspaceResponse, error)
}

func (g *grpcClient) CreateWorkspace(ctx context.Context, in *convert.CreateWorkspaceRequest) (*convert.CreateWorkspaceResponse, error) {

	protoResp, err := workspaceproto.NewWorkspaceServiceClient(g.conn).CreateWorkspace(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.CreateWorkspaceResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) DeleteWorkspace(ctx context.Context, in *convert.DeleteWorkspaceRequest) (*convert.DeleteWorkspaceResponse, error) {

	protoResp, err := workspaceproto.NewWorkspaceServiceClient(g.conn).DeleteWorkspace(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteWorkspaceResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) UpdateWorkspace(ctx context.Context, in *convert.UpdateWorkspaceRequest) (*convert.UpdateWorkspaceResponse, error) {

	protoResp, err := workspaceproto.NewWorkspaceServiceClient(g.conn).UpdateWorkspace(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.UpdateWorkspaceResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ListWorkspaces(ctx context.Context, in *convert.ListWorkspacesRequest) (*convert.ListWorkspacesResponse, error) {

	protoResp, err := workspaceproto.NewWorkspaceServiceClient(g.conn).ListWorkspace(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.ListWorkspacesResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) GetWorkspace(ctx context.Context, in *convert.GetWorkspaceRequest) (*convert.GetWorkspaceResponse, error) {

	protoResp, err := workspaceproto.NewWorkspaceServiceClient(g.conn).GetWorkspace(ctx, in.ToGRPC())
	if err != nil {
		return nil, err
	}
	out := &convert.GetWorkspaceResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (g *grpcClient) ImportWorkspace(ctx context.Context, in *convert.ImportWorkspaceRequest) (*convert.ImportWorkspaceResponse, error) {
	client := workspaceproto.NewWorkspaceServiceClient(g.conn)
	stream, err := client.ImportWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(in.FilePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// Maximum 1KB size per stream.
	buf := make([]byte, 1024)

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if err := stream.Send(&workspaceproto.ImportWorkspaceRequest{
			FileName: path.Base(in.FilePath),
			Content:  buf[:num],
			Storage: &workspaceproto.WorkspaceStorage{
				Nfs: &workspaceproto.NFSWorkspaceStorage{
					MountPath: in.MountPath,
				},
			},
		}); err != nil {
			return nil, err
		}
	}

	protoResp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	out := &convert.ImportWorkspaceResponse{}
	out.FromGRPC(protoResp)
	return out, nil
}

func (h *httpClient) CreateWorkspace(ctx context.Context, in *convert.CreateWorkspaceRequest) (*convert.CreateWorkspaceResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Post(h.url("workspace"))
	if err != nil {
		return nil, err
	}
	out := &convert.CreateWorkspaceResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) DeleteWorkspace(ctx context.Context, in *convert.DeleteWorkspaceRequest) (*convert.DeleteWorkspaceResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Delete(h.url("workspace/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.DeleteWorkspaceResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) UpdateWorkspace(ctx context.Context, in *convert.UpdateWorkspaceRequest) (*convert.UpdateWorkspaceResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Patch(h.url("workspace/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.UpdateWorkspaceResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) ListWorkspaces(ctx context.Context, in *convert.ListWorkspacesRequest) (*convert.ListWorkspacesResponse, error) {
	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace"))
	if err != nil {
		return nil, err
	}
	out := &convert.ListWorkspacesResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) ImportWorkspace(ctx context.Context, in *convert.ImportWorkspaceRequest) (*convert.ImportWorkspaceResponse, error) {
	req := h.rest.R().SetContext(ctx).ForceContentType("multipart/form-data").SetQueryParams(map[string]string{
		"mountPath": in.MountPath,
		"mountType": in.MountType,
	})
	file, err := os.Open(in.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	req.SetFile("file", path.Base(in.FilePath))
	httpResp, err := req.Put(h.url("workspace"))
	if err != nil {
		return nil, err
	}
	out := &convert.ImportWorkspaceResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}

func (h *httpClient) GetWorkspace(ctx context.Context, in *convert.GetWorkspaceRequest) (*convert.GetWorkspaceResponse, error) {

	req := h.restR(ctx)
	convert.AssignToHttpRequest(in, req)
	httpResp, err := req.Get(h.url("workspace/{id}"))
	if err != nil {
		return nil, err
	}
	out := &convert.GetWorkspaceResponse{}
	convert.AssignFromHttpResponse(httpResp, out)
	return out, nil
}
