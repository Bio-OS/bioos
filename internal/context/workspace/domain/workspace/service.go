package workspace

import (
	"context"
	"os"
	"path"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type Service interface {
	Import(ctx context.Context, workspaceID string, fileName string, storage Storage) error
}

type service struct {
	eventRepo  eventbus.EventRepository
	repository Repository
	eventbus   eventbus.EventBus
	factory    Factory
}

func (s *service) Import(ctx context.Context, workspaceID string, fileName string, storage Storage) error {
	baseDir := path.Join(storage.NFS.MountPath, workspaceID)
	zipFilePath := path.Join(baseDir, fileName)
	err := utils.Unzip(zipFilePath, baseDir)
	defer os.Remove(zipFilePath)
	if err != nil {
		applog.Errorw("unzip failed, clean all zipped files now", "err", err)
		os.RemoveAll(baseDir)
		return err
	}
	event := NewImportWorkspaceEvent(workspaceID, fileName, storage)
	return s.eventbus.Publish(ctx, event)
}

func NewService(repo Repository, eventRepo eventbus.EventRepository, bus eventbus.EventBus, factory Factory) Service {
	svc := &service{
		eventRepo:  eventRepo,
		repository: repo,
		eventbus:   bus,
		factory:    factory,
	}
	return svc
}
