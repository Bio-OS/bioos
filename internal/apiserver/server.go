package apiserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"k8s.io/component-base/featuregate"

	_ "github.com/Bio-OS/bioos/docs" // for swagger
	"github.com/Bio-OS/bioos/internal/apiserver/options"
	notebookserverapp "github.com/Bio-OS/bioos/internal/context/notebookserver/application"
	notebookservergrpc "github.com/Bio-OS/bioos/internal/context/notebookserver/interface/grpc"
	notebookserverproto "github.com/Bio-OS/bioos/internal/context/notebookserver/interface/grpc/proto"
	notebookserverhertz "github.com/Bio-OS/bioos/internal/context/notebookserver/interface/hertz"
	submissionapp "github.com/Bio-OS/bioos/internal/context/submission/application"
	submissiongrpc "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc"
	submissionproto "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc/proto"
	submissionhertz "github.com/Bio-OS/bioos/internal/context/submission/interface/hertz"
	workspaceapp "github.com/Bio-OS/bioos/internal/context/workspace/application"
	workspacegrpc "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	workspacehertz "github.com/Bio-OS/bioos/internal/context/workspace/interface/hertz"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/middlewares"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/server"
	"github.com/Bio-OS/bioos/pkg/version"
)

const (
	component = "bioos-apiserver"
)

func newBioosServerCommand(ctx context.Context, opts *options.Options) *cobra.Command {
	return &cobra.Command{
		Use:   component,
		Short: "bioos apiserver",
		Long: `bioos apiserver
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			if err := viper.Unmarshal(opts, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToTimeHookFunc(time.RFC3339),
				notebook.ResourceQuantityStringToInt64HookFunc,
			))); err != nil {
				return err
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			version.PrintVersionOrContinue()

			log.RegisterLogger(opts.LogOption)
			defer log.Sync()

			middlewares.RegisterAuthenticator(opts.AuthOption.AuthN)
			middlewares.RegisterAuthorizer(opts.AuthOption.AuthZ)

			log.Infow("server options", "options", opts)

			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Debugf("FLAG: --%s=%q", flag.Name, flag.Value)
			})

			return run(ctx, opts)
		},
		Args: cobra.ExactArgs(0),
	}
}

func run(ctx context.Context, opts *options.Options) error {
	log.Infow("called bioos-apiserver")

	workspaceService, err := workspaceapp.NewWorkspaceService(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		_ = workspaceService.Close(ctx)
	}()

	workspaceGRPCService := workspacegrpc.NewServer(workspaceService)
	// TODO invoke workspace API by grpc service client after solved backend token
	notebookserverService, err := notebookserverapp.NewService(ctx, opts, workspaceGRPCService)
	if err != nil {
		return fmt.Errorf("new notebook server service fail: %w", err)
	}
	defer func() {
		_ = notebookserverService.Close(ctx)
	}()

	submissionService, err := submissionapp.NewSubmissionService(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		_ = submissionService.Close(ctx)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", opts.ServerOption.Grpc.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	workflowGRPCService := workspacegrpc.NewWorkflowServer(workspaceService)
	datamodelGRPCService := workspacegrpc.NewDataModelServer(workspaceService)
	notebookGRPCService := workspacegrpc.NewNotebookServer(workspaceService)
	submissionGRPCService := submissiongrpc.NewSubmissionServer(submissionService)
	versionGRPCService := workspacegrpc.NewVersionServer()
	notebookserverGRPCService := notebookservergrpc.NewServer(notebookserverService)
	grpcServer, err := setupGrpcServer(
		opts,
		server.GetGRPCRegister(workspaceproto.RegisterWorkspaceServiceServer, workspaceGRPCService),
		server.GetGRPCRegister(workspaceproto.RegisterWorkflowServiceServer, workflowGRPCService),
		server.GetGRPCRegister(workspaceproto.RegisterDataModelServiceServer, datamodelGRPCService),
		server.GetGRPCRegister(workspaceproto.RegisterNotebookServiceServer, notebookGRPCService),
		server.GetGRPCRegister(submissionproto.RegisterSubmissionServiceServer, submissionGRPCService),
		server.GetGRPCRegister(workspaceproto.RegisterVersionServiceServer, versionGRPCService),
		server.GetGRPCRegister(notebookserverproto.RegisterNotebookServerServiceServer, notebookserverGRPCService),
	)
	if err != nil {
		log.Fatalf("failed to setup grpc server: %v", err)
	}
	httpServer, err := setupHTTPServer(
		opts,
		workspacehertz.NewRouteRegister(workspaceService),
		submissionhertz.NewRouteRegister(submissionService),
		notebookserverhertz.NewRouteRegister(notebookserverService),
	)
	if err != nil {
		log.Fatalf("failed to setup http server: %v", err)
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		return err
	})
	eg.Go(func() error {
		httpServer.Spin()
		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalf("get error: %v", err)
	}

	return nil
}

// NewBioosServerCommand instance a bioos server command.
func NewBioosServerCommand(ctx context.Context) *cobra.Command {
	opt := options.NewOptions()
	cmd := newBioosServerCommand(ctx, opt)

	opt.AddFlags(cmd.Flags())
	featureGate := featuregate.NewFeatureGate()
	featureGate.AddFlag(cmd.Flags())
	version.AddFlags(cmd.Flags())
	cmd.Flags().AddFlag(pflag.Lookup(options.ConfigFlagName))
	return cmd
}
