package eventbus

import (
	"context"
	"time"
)

type Event struct {
	EventID     string
	Type        string
	Payload     string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ScheduledAt time.Time
	RetryCount  int // The number of times the events has been retried
}

type Filter struct {
	Type    []string
	Payload string
	Status  []string
}

const (
	EventStatusPending   = "pending"
	EventStatusDequeue   = "dequeue"
	EventStatusRunning   = "running"
	EventStatusCompleted = "completed"
	EventStatusFailed    = "failed"
)

// ErrEventRunningDelayed stand for long-run event that need to delay add back.
type ErrEventRunningDelayed struct {
	message string
	delay   time.Duration
}

func NewErrEventRunningDelayed(message string, delay time.Duration) ErrEventRunningDelayed {
	return ErrEventRunningDelayed{
		message: message,
		delay:   delay,
	}
}

func (e ErrEventRunningDelayed) Error() string {
	return e.message
}

func (e ErrEventRunningDelayed) Delay() time.Duration {
	return e.delay
}

// EventHandler represents a function that handle a event
type EventHandler interface {
	Handle(ctx context.Context, payload string) error
}
type EventHandlerFunc func(ctx context.Context, payload string) error

// Handle calls f(ctx, event).
func (f EventHandlerFunc) Handle(ctx context.Context, payload string) error {
	return f(ctx, payload)
}

type EventRepository interface {
	Get(ctx context.Context, id string) (*Event, error)
	Save(ctx context.Context, event *Event) error
	ListAndLockUnfinishedEvents(ctx context.Context, limit int, eventTypes []string) ([]*Event, error)
	UpdateStatus(ctx context.Context, event *Event, status string) error
	UpdateRetryCount(ctx context.Context, event *Event, retryCount int) error
	Search(ctx context.Context, filter *Filter) ([]*Event, error)
}
