package command

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
)

type SwitchCommand struct {
	ID          string
	WorkspaceID string
	OnOff       bool // true is trun on, false is turn off
}

type SwitchHandler interface {
	Handle(context.Context, *SwitchCommand) error
}

func NewSwitchHandler(svc domain.Service) SwitchHandler {
	return &switchHandler{
		service: svc,
	}
}

type switchHandler struct {
	service domain.Service
}

func (h *switchHandler) Handle(ctx context.Context, cmd *SwitchCommand) error {
	// TODO check workspace exist
	var err error
	if cmd.OnOff {
		err = h.service.Start(ctx, cmd.ID)
	} else {
		err = h.service.Stop(ctx, cmd.ID)
	}
	return err
}
