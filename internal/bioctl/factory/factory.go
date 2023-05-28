package factory

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"

	clioptions "github.com/Bio-OS/bioos/internal/bioctl/options"
	"github.com/Bio-OS/bioos/internal/bioctl/utils"
	"github.com/Bio-OS/bioos/pkg/client"
	pkgutils "github.com/Bio-OS/bioos/pkg/utils"
)

type Factory interface {
	WorkspaceClient() (WorkspaceClient, error)
	DataModelClient() (DataModelClient, error)
	WorkflowClient() (WorkflowClient, error)
	VersionClient() (VersionClient, error)
	SubmissionClient() (SubmissionClient, error)
}

func NewFactory(opts *clioptions.ClientOptions) Factory {
	return factoryImpl{
		opts: opts,
	}
}

var _ Factory = factoryImpl{}

type factoryImpl struct {
	opts *clioptions.ClientOptions
}

func (f factoryImpl) newGrpcClient() (*grpcClient, error) {
	conn, err := pkgutils.GrpcDial(f.opts.ConnectInfo, f.opts.AuthInfo)
	if err != nil {
		utils.CheckErr(err)
	}
	return &grpcClient{
		conn: conn,
	}, nil
}

func (f factoryImpl) newHttpClient() (*httpClient, error) {
	transport, err := pkgutils.NewTransportWithAuth(f.opts.ConnectInfo, f.opts.AuthInfo)
	if err != nil {
		return nil, err
	}
	serverAddr := f.opts.ServerAddr
	if f.opts.Insecure {
		serverAddr = "http://" + serverAddr
	} else {
		serverAddr = "https://" + serverAddr
	}
	return &httpClient{
		addr: strings.TrimRight(serverAddr, "/"),
		rest: resty.NewWithClient(transport.Client()),
	}, err
}

func (f factoryImpl) WorkspaceClient() (WorkspaceClient, error) {
	if err := f.opts.Method.Validate(); err != nil {
		return nil, err
	}
	switch f.opts.Method {
	case client.GRPCMethod:
		return f.newGrpcClient()
	case client.HTTPMethod:
		return f.newHttpClient()
	}
	return nil, nil
}

func (f factoryImpl) DataModelClient() (DataModelClient, error) {
	if err := f.opts.Method.Validate(); err != nil {
		return nil, err
	}
	switch f.opts.Method {
	case client.GRPCMethod:
		return f.newGrpcClient()
	case client.HTTPMethod:
		return f.newHttpClient()
	}
	return nil, nil
}

func (f factoryImpl) VersionClient() (VersionClient, error) {
	if err := f.opts.Method.Validate(); err != nil {
		return nil, err
	}
	switch f.opts.Method {
	case client.GRPCMethod:
		return f.newGrpcClient()
	case client.HTTPMethod:
		return f.newHttpClient()
	}
	return nil, nil
}

func (f factoryImpl) SubmissionClient() (SubmissionClient, error) {
	if err := f.opts.Method.Validate(); err != nil {
		return nil, err
	}
	switch f.opts.Method {
	case client.GRPCMethod:
		return f.newGrpcClient()
	case client.HTTPMethod:
		return f.newHttpClient()
	}
	return nil, nil
}

func (f factoryImpl) WorkflowClient() (WorkflowClient, error) {
	if err := f.opts.Method.Validate(); err != nil {
		return nil, err
	}
	switch f.opts.Method {
	case client.GRPCMethod:
		return f.newGrpcClient()
	case client.HTTPMethod:
		return f.newHttpClient()
	}
	return nil, nil
}

type grpcClient struct {
	conn *grpc.ClientConn
}

type httpClient struct {
	addr string
	rest *resty.Client
}

func (h *httpClient) restR(ctx context.Context) *resty.Request {
	return h.rest.R().SetContext(ctx).ForceContentType("application/json")
}

func (h *httpClient) url(apipath string, pathParams ...interface{}) string {
	return strings.Join([]string{h.addr, fmt.Sprintf(apipath, pathParams...)}, "/")
}
