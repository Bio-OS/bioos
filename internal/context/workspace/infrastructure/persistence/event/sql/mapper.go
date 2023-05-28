package sql

import (
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

func eventPOToEventDO(e *Event) *eventbus.Event {
	return &eventbus.Event{
		EventID:     e.EventID,
		Type:        e.Type,
		Payload:     e.Payload,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		ScheduledAt: e.ScheduledAt,
		RetryCount:  e.RetryCount,
	}

}

func eventDOToEventPO(e *eventbus.Event) *Event {
	return &Event{
		EventID:     e.EventID,
		Type:        e.Type,
		Payload:     e.Payload,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		ScheduledAt: e.ScheduledAt,
		RetryCount:  e.RetryCount,
	}
}
