package apiserver

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/Bio-OS/bioos/internal/apiserver/options"
	middlewaregrpc "github.com/Bio-OS/bioos/pkg/middlewares/grpc"
	appserver "github.com/Bio-OS/bioos/pkg/server"
)

func setupGrpcServer(opts *options.Options, registers ...appserver.GRPCRegister) (*grpc.Server, error) {
	serverOptions := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			middlewaregrpc.NewAuthStreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
			middlewaregrpc.RBACStreamServerChain(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			middlewaregrpc.NewAuthUnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			middlewaregrpc.RBACUnaryServerChain(),
		)),
	}
	if opts.ServerOption.Grpc.TLS {
		tlsCredentials, err := loadTLSCredentials(opts.ServerOption.CertFile, opts.ServerOption.KeyFile, opts.ServerOption.CaFile)
		if err != nil {
			return nil, err
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}
	grpcServer := grpc.NewServer(serverOptions...)
	for _, r := range registers {
		r(grpcServer)
	}
	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)
	hs := health.NewServer()
	hs.SetServingStatus("grpc.health.v1.workspaceservice", 1)
	healthpb.RegisterHealthServer(grpcServer, hs)
	return grpcServer, nil
}

func loadTLSCredentials(certFile, keyFile, caFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(ca)
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    certPool,
	}), nil
}
