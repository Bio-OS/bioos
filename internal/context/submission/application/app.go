package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/internal/apiserver/options"
	runcommand "github.com/Bio-OS/bioos/internal/context/submission/application/command/run"
	submissioncommand "github.com/Bio-OS/bioos/internal/context/submission/application/command/submission"
	runquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	submissionquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/run"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	runsqlpo "github.com/Bio-OS/bioos/internal/context/submission/infrastructure/persistence/run/sql"
	submissionsqlpo "github.com/Bio-OS/bioos/internal/context/submission/infrastructure/persistence/submission/sql"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	eventmongopo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/event/mongo"
	eventsqlpo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/event/sql"
	"github.com/Bio-OS/bioos/pkg/client"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type closer func(ctx context.Context) error

type SubmissionService struct {
	SubmissionCommands *submissioncommand.Commands
	SubmissionQueries  *submissionquery.Queries
	RunCommands        *runcommand.Commands
	RunQueries         *runquery.Queries
	closer             closer
}

func (w *SubmissionService) Close(ctx context.Context) error {
	if w.closer != nil {
		if err := w.closer(ctx); err != nil {
			return err
		}
	}

	return nil
}

func NewSubmissionService(ctx context.Context, opts *options.Options) (*SubmissionService, error) {
	var (
		err                 error
		dbCloser            closer
		submissionRepo      submission.Repository
		submissionReadModel submissionquery.ReadModel
		runRepo             run.Repository
		runReadModel        runquery.ReadModel
		eventRepo           eventbus.EventRepository
		eventBus            eventbus.EventBus
		grpcFactory         grpc.Factory
		wesClient           wes.Client
	)

	if opts.DBOption.Mongo != nil && opts.DBOption.Mongo.Enabled() {
		_, mongoDB, err := opts.DBOption.Mongo.GetDBInstance(ctx)
		if err != nil {
			return nil, fmt.Errorf("get mongo db client fail: %w", err)
		}
		// todo add submission&run mongo repository and read model
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
		} else if opts.DBOption.SQLite3 != nil || opts.DBOption.SQLite3.Enabled() {
			if orm, err = opts.DBOption.SQLite3.GetGORMInstance(ctx); err != nil {
				return nil, fmt.Errorf("get mysql db client fail: %w", err)
			}
		} else {
			return nil, fmt.Errorf("none sql db options")
		}
		if submissionRepo, err = submissionsqlpo.NewSubmissionRepository(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if submissionReadModel, err = submissionsqlpo.NewSubmissionReadModel(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql read model fail: %w", err)
		}
		if runRepo, err = runsqlpo.NewRunRepository(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if runReadModel, err = runsqlpo.NewRunReadModel(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql read model fail: %w", err)
		}
		if eventRepo, err = eventsqlpo.NewEventRepository(ctx, orm, opts.EventBusOption.DequeueTimeout, opts.EventBusOption.RunningTimeout); err != nil {
			return nil, fmt.Errorf("new sql event repository fail: %w", err)
		}
	}

	if opts.Client.Method == client.GRPCMethod {
		grpcFactory = grpc.NewFactory(opts.Client)
	}

	wesClient = wes.NewClient(opts.WesOption)

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
			log.Errorw("start event bus failed", "err", err)
		}
	}()

	submissionFactory := submission.NewSubmissionFactory(ctx)
	return &SubmissionService{
		SubmissionCommands: submissioncommand.NewCommands(grpcFactory, submissionRepo, submissionFactory, eventBus, submissionReadModel, runReadModel),
		SubmissionQueries:  submissionquery.NewQueries(grpcFactory, submissionReadModel),
		RunCommands:        runcommand.NewCommands(grpcFactory, runRepo, eventBus, submissionReadModel, wesClient),
		RunQueries:         runquery.NewQueries(grpcFactory, runReadModel, submissionReadModel),
		closer:             dbCloser,
	}, nil
}
