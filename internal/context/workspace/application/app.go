package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/internal/apiserver/options"
	datamodelcommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/data-model"
	notebookcommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/notebook"
	workflowcommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workflow"
	workspacecommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workspace"
	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	notebookquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	workflowquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	datamodelsqlpo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/data-model/sql"
	eventmongopo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/event/mongo"
	eventsqlpo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/event/sql"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/notebook/filesystem"
	workflowmongo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/workflow/mongo"
	workflowsql "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/workflow/sql"
	workspacemongo "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/workspace/mongo"
	workspacesql "github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/persistence/workspace/sql"
	"github.com/Bio-OS/bioos/pkg/log"
)

type closer func(ctx context.Context) error

type WorkspaceService struct {
	WorkspaceCommands *workspacecommand.Commands
	WorkspaceQueries  *workspacequery.Queries
	NotebookCommands  *notebookcommand.Commands
	NotebookQueries   *notebookquery.Queries
	DataModelCommands *datamodelcommand.Commands
	DataModelQueries  *datamodelquery.Queries
	WorkflowCommands  *workflowcommand.Commands
	WorkflowQueries   *workflowquery.Queries

	closer closer
}

func (w *WorkspaceService) Close(ctx context.Context) error {
	if w.closer != nil {
		if err := w.closer(ctx); err != nil {
			return err
		}
	}

	return nil
}

func NewWorkspaceService(ctx context.Context, opts *options.Options) (*WorkspaceService, error) {
	var (
		err                error
		dbCloser           closer
		workspaceRepo      workspace.Repository
		workspaceReadModel workspacequery.WorkspaceReadModel
		dataModelRepo      datamodel.Repository
		dataModelReadModel datamodelquery.DataModelReadModel
		workflowRepo       workflow.Repository
		workflowReadModel  workflowquery.ReadModel
		eventRepo          eventbus.EventRepository
		eventBus           eventbus.EventBus
	)

	if opts.DBOption.Mongo != nil && opts.DBOption.Mongo.Enabled() {
		mongoClient, mongoDB, err := opts.DBOption.Mongo.GetDBInstance(ctx)
		if err != nil {
			return nil, fmt.Errorf("get mongo db client fail: %w", err)
		}
		if workspaceRepo, err = workspacemongo.NewWorkspaceRepository(ctx, mongoDB); err != nil {
			return nil, fmt.Errorf("new mongodb repository fail: %w", err)
		}
		if workspaceReadModel, err = workspacemongo.NewWorkspaceReadModel(ctx, mongoDB); err != nil {
			return nil, fmt.Errorf("new mongodb read model fail: %w", err)
		}
		if workflowRepo, err = workflowmongo.NewRepository(ctx, mongoDB, mongoClient); err != nil {
			return nil, fmt.Errorf("new mongodb repository fail: %w", err)
		}
		if workflowReadModel, err = workflowmongo.NewReadModel(ctx, mongoDB); err != nil {
			return nil, fmt.Errorf("new mongodb read model fail: %w", err)
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
		} else if opts.DBOption.SQLite3 != nil || opts.DBOption.SQLite3.Enabled() {
			if orm, err = opts.DBOption.SQLite3.GetGORMInstance(ctx); err != nil {
				return nil, fmt.Errorf("get mysql db client fail: %w", err)
			}
		} else {
			return nil, fmt.Errorf("none sql db options")
		}
		if workspaceRepo, err = workspacesql.NewWorkspaceRepository(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if workspaceReadModel, err = workspacesql.NewWorkspaceReadModel(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql read model fail: %w", err)
		}
		if workflowRepo, err = workflowsql.NewRepository(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if workflowReadModel, err = workflowsql.NewReadModel(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql read model fail: %w", err)
		}
		if dataModelRepo, err = datamodelsqlpo.NewDataModelRepository(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql repository fail: %w", err)
		}
		if dataModelReadModel, err = datamodelsqlpo.NewDataModelReadModel(ctx, orm); err != nil {
			return nil, fmt.Errorf("new sql read model fail: %w", err)
		}
		if eventRepo, err = eventsqlpo.NewEventRepository(ctx, orm, opts.EventBusOption.DequeueTimeout, opts.EventBusOption.RunningTimeout); err != nil {
			return nil, fmt.Errorf("new sql event repository fail: %w", err)
		}
	}
	eOpts := []eventbus.Option{
		eventbus.WithMaxRetries(opts.EventBusOption.MaxRetries),
		eventbus.WithSyncPeriod(opts.EventBusOption.SyncPeriod),
		eventbus.WithBatchSize(opts.EventBusOption.BatchSize),
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

	var notebookRepo notebook.Repository
	var notebookReadModel notebookquery.ReadModel
	if opts.StorageOption.FileSystem != nil && opts.StorageOption.FileSystem.Enabled() {
		if notebookRepo, err = filesystem.NewRepository(opts.StorageOption.FileSystem.NotebookRootPath()); err != nil {
			return nil, fmt.Errorf("new notebook fs repository fail: %w", err)
		}
		if notebookReadModel, err = filesystem.NewReadModel(opts.StorageOption.FileSystem.NotebookRootPath()); err != nil {
			return nil, fmt.Errorf("new notebook fs read model fail: %w", err)
		}
	} else {
		return nil, fmt.Errorf("none storage options")
	}

	workspaceFactory := workspace.NewWorkspaceFactory(ctx)
	dataModelFactory := datamodel.NewDataModelFactory()
	workflowFactory := workflow.NewFactory(ctx)
	notebookFactory := notebook.NewFactory()

	return &WorkspaceService{
		WorkspaceCommands: workspacecommand.NewCommands(workspaceRepo, eventRepo, workspaceFactory, eventBus),
		WorkspaceQueries:  workspacequery.NewQueries(workspaceReadModel),
		WorkflowCommands:  workflowcommand.NewCommands(workflowRepo, workflowReadModel, workflowFactory, workspaceReadModel, eventBus, opts.ServerOption.WomtoolFile),
		WorkflowQueries:   workflowquery.NewQueries(workflowReadModel, workspaceReadModel),
		NotebookCommands:  notebookcommand.NewCommands(notebookRepo, workspaceReadModel, eventBus, notebookReadModel, notebookFactory),
		NotebookQueries:   notebookquery.NewQueries(notebookReadModel, workspaceReadModel),
		DataModelCommands: datamodelcommand.NewCommands(dataModelRepo, workspaceReadModel, dataModelFactory, dataModelReadModel, eventBus),
		DataModelQueries:  datamodelquery.NewQueries(workspaceReadModel, dataModelReadModel),
		closer:            dbCloser,
	}, nil
}
