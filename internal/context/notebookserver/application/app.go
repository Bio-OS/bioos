package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/internal/apiserver/options"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/command"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/infrastructure/k8shub"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/infrastructure/persistence/mongo"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/infrastructure/persistence/sql"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	eventmongopo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/event/mongo"
	eventsqlpo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/event/sql"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type closer func(ctx context.Context) error

type Service struct {
	Commands *command.Commands
	Queries  *query.Queries
	closer   closer
}

func (s *Service) Close(ctx context.Context) error {
	if s.closer != nil {
		if err := s.closer(ctx); err != nil {
			return err
		}
	}

	return nil
}

func NewService(ctx context.Context, opts *options.Options, workspaceService proto.WorkspaceServiceServer) (*Service, error) {
	var (
		err       error
		dbCloser  closer
		repo      domain.Repository
		readModel query.ReadModel
		eventRepo eventbus.EventRepository
		eventBus  eventbus.EventBus
	)

	// init database
	if opts.DBOption.Mongo != nil && opts.DBOption.Mongo.Enabled() {
		_, mongoDB, err := opts.DBOption.Mongo.GetDBInstance(ctx)
		if err != nil {
			return nil, fmt.Errorf("get mongo db client fail: %w", err)
		}
		if repo, err = mongo.NewRepository(ctx, mongoDB); err != nil {
			return nil, fmt.Errorf("new mongodb repository fail: %w", err)
		}
		if readModel, err = mongo.NewReadModel(ctx, mongoDB); err != nil {
			return nil, fmt.Errorf("new mongodb repository fail: %w", err)
		}
		if eventRepo, err = eventmongopo.NewEventRepository(ctx, mongoDB, opts.EventBusOption.DequeueTimeout, opts.EventBusOption.RunningTimeout); err != nil {
			return nil, fmt.Errorf("new mongodb event repository fail: %w", err)
		}
		dbCloser = func(ctx context.Context) error {
			return mongoDB.Client().Disconnect(ctx)
		}
	} else {
		// deal various sql db with gorm
		var orm *gorm.DB
		if opts.DBOption.MySQL != nil && opts.DBOption.MySQL.Enabled() {
			if orm, err = opts.DBOption.MySQL.GetGORMInstance(ctx); err != nil {
				return nil, fmt.Errorf("get mysql db client fail: %w", err)
			}
			dbCloser = func(ctx context.Context) error {
				dbInstance, _ := orm.DB()
				return dbInstance.Close()
			}
		} else if opts.DBOption.SQLite3 != nil && opts.DBOption.SQLite3.Enabled() {
			if orm, err = opts.DBOption.SQLite3.GetGORMInstance(ctx); err != nil {
				return nil, fmt.Errorf("get mysql db client fail: %w", err)
			}
		} else {
			return nil, fmt.Errorf("none sql db options")
		}
		if repo, err = sql.NewRepository(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if readModel, err = sql.NewReadModel(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if eventRepo, err = eventsqlpo.NewEventRepository(ctx, orm, opts.EventBusOption.DequeueTimeout, opts.EventBusOption.RunningTimeout); err != nil {
			return nil, fmt.Errorf("new sql event repository fail: %w", err)
		}
	}

	eOpts := []eventbus.Option{
		eventbus.WithMaxRetries(opts.EventBusOption.MaxRetries),
		eventbus.WithSyncPeriod(opts.EventBusOption.SyncPeriod),
		eventbus.WithBatchSize(opts.EventBusOption.BatchSize),
	}
	if eventBus, err = eventbus.NewEventBus(eventRepo, eOpts...); err != nil {
		return nil, err
	}
	go func() {
		err := eventBus.Start(ctx, opts.EventBusOption.Workers)
		if err != nil {
			applog.Errorw("start event bus failed", "err", err)
		}
	}()

	// generate jupyter runtime
	var runtime domain.Runtime
	if opts.NotebookOption.StaticJupyterhub.Endpoint != "" {
		runtime, err = k8shub.NewRuntime(
			ctx,
			&opts.NotebookOption.StaticJupyterhub,
			opts.StorageOption,
		)
		if err != nil {
			return nil, fmt.Errorf("can not new k8s jupyterhub runtime: %w", err)
		}
	} else {
		runtime = domain.UnimplementedRuntime{}
	}

	factory := domain.NewFactory(
		opts.NotebookOption.ListOfficialImages(),
		opts.NotebookOption.ResourceSizes,
	)
	return &Service{
		Commands: command.NewCommands(repo, factory, runtime, workspaceService, opts.StorageOption, eventBus),
		Queries:  query.NewQueries(readModel, runtime),
		closer:   dbCloser,
	}, nil
}
