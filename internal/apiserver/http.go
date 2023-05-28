package apiserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/google/uuid"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/http2/factory"
	"github.com/hertz-contrib/requestid"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"golang.org/x/net/http2"

	"github.com/Bio-OS/bioos/internal/apiserver/options"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	apphertz "github.com/Bio-OS/bioos/pkg/middlewares/hertz"
	"github.com/Bio-OS/bioos/pkg/notebook"
	appserver "github.com/Bio-OS/bioos/pkg/server"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/version"
)

func setupHTTPServer(opts *options.Options, registers ...appserver.RouteRegister) (*server.Hertz, error) {
	serverOptions := []config.Option{
		server.WithHostPorts(fmt.Sprintf(":%s", opts.ServerOption.Http.Port)),
		server.WithMaxRequestBodySize(opts.ServerOption.Http.MaxRequestBodySize),
		server.WithALPN(true),
	}

	if opts.ServerOption.Http.TLS {
		cert, err := tls.LoadX509KeyPair(opts.ServerOption.CertFile, opts.ServerOption.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("fail to load x509 pair: %w", err)
		}
		certBytes, err := os.ReadFile(opts.ServerOption.CaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read ca file: %w", err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			return nil, fmt.Errorf("failed to parse ca file: %w", err)
		}
		// set server tls.Config
		cfg := &tls.Config{
			// add certificate
			Certificates: []tls.Certificate{cert},
			MaxVersion:   tls.VersionTLS13,
			// enable client authentication
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
			// cipher suites supported
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
			// set application protocol http2
			NextProtos: []string{http2.NextProtoTLS},
		}
		serverOptions = append(serverOptions, server.WithTLS(cfg))
	}
	httpServer := server.Default(serverOptions...)

	httpServer.AddProtocol(http2.NextProtoTLS, factory.NewServerFactory())
	setupMiddlewares(httpServer)
	setupRouter(httpServer, opts)
	for _, r := range registers {
		r.AddRoute(httpServer)
	}
	return httpServer, nil
}

func setupMiddlewares(h *server.Hertz) {
	h.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "GET", "POST", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowWebSockets:  true,
			AllowWildcard:    true,
			MaxAge:           12 * time.Hour,
		}),
		requestid.New(
			requestid.WithGenerator(func(ctx context.Context, c *app.RequestContext) string {
				return uuid.New().String()
			}),
			// set custom header for request id
			requestid.WithCustomHeaderStrKey(consts.XRequestIDKey),
		),
		apphertz.Logger(),
	)
}

func setupRouter(h *server.Hertz, opts *options.Options) {
	h.GET("/ping", PingHandler)
	h.GET("/version", VersionHandler)
	url := swagger.URL("/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	h.GET("/.well-known/configuration", clientConfigHandler(opts))

	h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		utils.WriteHertzErrorResponse(c, &apperrors.AppError{Code: apperrors.RouteNotFoundCode, Message: "route not found"})
	})
}

type clientConfig struct {
	Storage struct {
		FSPath []string `json:"fsPath,omitempty"`
	} `json:"storage"`
	Notebook struct {
		OfficialImages  []notebook.Image        `json:"officialImages"`
		ResourceOptions []notebook.ResourceSize `json:"resourceOptions"`
	} `json:"notebook"`
}

// GetClientConfig get client configuration
//
//	@Summary		use to get client configuration
//	@Description	get client configuration
//	@Router			/.well-known/configuration [get]
//	@Produce		application/json
//	@Success		200	{object}	clientConfig
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func clientConfigHandler(opts *options.Options) app.HandlerFunc {
	return func(_ context.Context, ctx *app.RequestContext) {
		var resp clientConfig
		if opts.StorageOption != nil {
			if opts.StorageOption.FileSystem != nil {
				// string array is for multi share storage porpose in future
				resp.Storage.FSPath = []string{opts.StorageOption.FileSystem.RootPath}
			}
		}
		resp.Notebook.OfficialImages = opts.NotebookOption.OfficialImages
		resp.Notebook.ResourceOptions = opts.NotebookOption.ListResourceSizes()
		utils.WriteHertzOKResponse(ctx, &resp)
	}
}

// PingHandler ping handler
//
//	@Summary		ping
//	@Description	ping
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/ping [get]
//	@Success		200
func PingHandler(_ context.Context, ctx *app.RequestContext) {
	ctx.JSON(http.StatusOK, map[string]string{
		"ping": "pong",
	})
}

// VersionHandler version handler
//
//	@Summary		version Summary
//	@Description	version Description
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/version [get]
//	@Success		200
func VersionHandler(_ context.Context, ctx *app.RequestContext) {
	ctx.JSON(http.StatusOK, version.Get())
}
