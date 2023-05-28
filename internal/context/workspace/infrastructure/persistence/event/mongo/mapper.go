package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

func eventPOToEventDO(e *Event) *eventbus.Event {
	return &eventbus.Event{
		EventID:     e.EventID,
		Type:        e.Type,
		Payload:     e.Payload,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt.Time(),
		UpdatedAt:   e.UpdatedAt.Time(),
		ScheduledAt: e.ScheduledAt.Time(),
		RetryCount:  e.RetryCount,
	}

}

func eventDOToEventPO(e *eventbus.Event) *Event {
	return &Event{
		EventID:     e.EventID,
		Type:        e.Type,
		Payload:     e.Payload,
		Status:      e.Status,
		CreatedAt:   primitive.NewDateTimeFromTime(e.CreatedAt),
		UpdatedAt:   primitive.NewDateTimeFromTime(e.UpdatedAt),
		ScheduledAt: primitive.NewDateTimeFromTime(e.ScheduledAt),
		RetryCount:  e.RetryCount,
	}
}
