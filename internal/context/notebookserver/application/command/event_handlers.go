package command

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/storage"
)

func addEventHandle(eb eventbus.EventBus, svc domain.Service, factory *domain.Factory, workspaceClient proto.WorkspaceServiceServer, storageOpts *storage.Options) {
	eb.Subscribe(domain.ImportNotebookServers, &importNotebookServersHandler{
		service:         svc,
		eventbus:        eb,
		factory:         factory,
		workspaceClient: workspaceClient,
		storageOpts:     storageOpts,
	})
}

type importNotebookServersHandler struct {
	service         domain.Service
	workspaceClient proto.WorkspaceServiceServer
	storageOpts     *storage.Options
	eventbus        eventbus.EventBus
	factory         *domain.Factory
}

func (h *importNotebookServersHandler) Handle(ctx context.Context, payload string) error {
	log.Infow("start to consume import notebook servers event", "payload", payload)
	event, err := domain.NewImportNotebookServersEventFromPayload([]byte(payload))
	if err != nil {
		return fmt.Errorf("decode event payload fail: %w", err)
	}

	handler := NewCreateHandler(h.service, h.factory, h.workspaceClient, h.storageOpts)

	var defaultResourceSize notebook.ResourceSize
	registeredSizes := h.factory.GetResourceSizes()
	if len(registeredSizes) > 0 {
		defaultResourceSize = registeredSizes[0]
	} else {
		return fmt.Errorf("notebookservers factory has no registered resourceSize")
	}

	_, err = handler.Handle(ctx, &CreateCommand{
		WorkspaceID:  event.WorkspaceID,
		Image:        event.Schema.Image.Name,
		ResourceSize: defaultResourceSize,
	})
	if err != nil {
		return err
	}
	return nil
}
