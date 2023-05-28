package eventbus

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"

	applog "github.com/Bio-OS/bioos/pkg/log"
)

// EventBus stands for event bus.
type EventBus interface {
	Publish(ctx context.Context, event IEvent) error
	Subscribe(eventType string, handler EventHandler)
	Start(ctx context.Context, workers int) error
	Close(ctx context.Context) error
}

type IEvent interface {
	EventType() string
	Payload() []byte
	Delay() time.Duration
}

// Impl implement event bus
type Impl struct {
	sync.Mutex
	repository  EventRepository
	subscribers map[string][]EventHandler
	maxRetries  int
	syncPeriod  time.Duration
	batchSize   int
	queue       workqueue.RateLimitingInterface
	runningSet  sets.Set[string]
}

var _ EventBus = &Impl{}

// NewEventBus new an event bus.
func NewEventBus(repository EventRepository, options ...Option) (EventBus, error) {
	impl := &Impl{
		repository:  repository,
		subscribers: make(map[string][]EventHandler),
		queue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "event-bus"),
		runningSet:  sets.New[string](),
	}
	for _, option := range options {
		option(impl)
	}
	return impl, nil
}

// Publish publish an event with payload and delay.
func (engine *Impl) Publish(ctx context.Context, iEvent IEvent) error {
	now := time.Now()
	event := &Event{
		EventID:     uuid.New().String(),
		Type:        iEvent.EventType(),
		Payload:     string(iEvent.Payload()),
		Status:      EventStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
		ScheduledAt: now,
	}
	if iEvent.Delay() > 0 {
		event.ScheduledAt = now.Add(iEvent.Delay())
	}
	// TODO need to gc events in repository
	return retry.OnError(retry.DefaultRetry, func(err error) bool {
		return err != nil
	}, func() error {
		return engine.repository.Save(ctx, event)
	})
}

// Subscribe register a handler for the event type.
func (engine *Impl) Subscribe(eventType string, handler EventHandler) {
	engine.Lock()
	defer engine.Unlock()
	engine.subscribers[eventType] = append(engine.subscribers[eventType], handler)
}

// Start start event bus
func (engine *Impl) Start(ctx context.Context, workers int) error {
	// don't let panics crash the process
	runtime.ReallyCrash = false
	defer runtime.HandleCrash()
	// make sure the work queue is shutdown which will trigger workers to end
	defer engine.queue.ShutDown()

	for i := 0; i < workers; i++ {
		go wait.Until(func() { engine.runWorker(ctx) }, time.Second, ctx.Done())
	}
	engine.processPendingEvents(ctx)
	return nil
}

// Close exit event bus
func (engine *Impl) Close(ctx context.Context) error {
	return nil
}

func (engine *Impl) processEvent(ctx context.Context, event *Event) error {
	errs := make([]error, 0)
	handlers := engine.subscribers[event.Type]
	if len(handlers) == 0 { // no handlers just return
		applog.Infow("no handlers for event", "eventType", event.Type)
		return nil
	}

	// Check if the event is scheduled and if its scheduled time has passed before executing it
	if !event.ScheduledAt.IsZero() && event.ScheduledAt.After(time.Now()) {
		delay := event.ScheduledAt.Sub(time.Now())
		engine.queue.AddAfter(event.EventID, delay)
		return nil
	}

	// if ready reach max retry times, mark it to failed
	if event.RetryCount >= engine.maxRetries {
		return engine.repository.UpdateStatus(ctx, event, EventStatusFailed)
	}
	// if it is already running, but can take it, have two circumstance:
	// 1) short task and outdated
	// 2) long time task
	if event.Status == EventStatusRunning && !engine.runningSet.Has(event.EventID) {
		return engine.repository.UpdateStatus(ctx, event, EventStatusFailed)
	}

	// mark event running to prevent handle concurrently
	if marked := engine.markEventRunning(event); !marked { // event is running, skip
		return nil
	}
	defer engine.unmarkEventRunning(event)
	// Set the status of the event to "running" in db
	if err := engine.repository.UpdateStatus(ctx, event, EventStatusRunning); err != nil {
		return err
	}

	var wg sync.WaitGroup
	runningFlag := false
	for _, h := range handlers {
		wg.Add(1)
		handler := h
		func() { // 同一个任务多个handler不要并发 避免锁
			defer wg.Done()
			if err := handler.Handle(ctx, event.Payload); err != nil {
				if delayedErr, ok := err.(ErrEventRunningDelayed); ok {
					// If the event is still running and needs to be delayed, schedule it for execution after the specified delay
					engine.queue.AddAfter(event.EventID, delayedErr.Delay())
					runningFlag = true
				} else {
					errs = append(errs, err)
				}
			}
		}()
	}
	wg.Wait()

	// Update the status of the event based on the result of the handler
	if len(errs) > 0 {
		event.RetryCount++ // Increment the retry count
		if err := engine.repository.UpdateRetryCount(ctx, event, event.RetryCount); err != nil {
			errs = append(errs, err)
		}
		// if ready reach max retry times, mark it to failed
		if event.RetryCount >= engine.maxRetries {
			return engine.repository.UpdateStatus(ctx, event, EventStatusFailed)
		}
		// else update task status to pending let it retry
		if err := engine.repository.UpdateStatus(ctx, event, EventStatusPending); err != nil {
			return err
		}
		return errors.NewAggregate(errs)
	} else {
		// avoid change running event to completed
		if !runningFlag {
			if err := engine.repository.UpdateStatus(ctx, event, EventStatusCompleted); err != nil {
				return err
			}
		}

		return nil
	}
}

func (engine *Impl) processPendingEvents(ctx context.Context) {
	ticker := time.Tick(engine.syncPeriod)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker: // Check for scheduled events every minute
			// if still have tasks not consumed, not fetch new task
			if engine.queue.Len() > 0 {
				continue
			}

			// only fetch those events with handler
			if len(engine.subscribers) > 0 {
				events, err := engine.repository.ListAndLockUnfinishedEvents(ctx, engine.batchSize, maps.Keys(engine.subscribers))
				if err != nil {
					applog.Errorw("Error loading scheduled events", "err", err)
				} else {
					for _, event := range events {
						engine.queue.Add(event.EventID)
					}
				}
			}
		}
	}
}

func (engine *Impl) runWorker(ctx context.Context) {
	for engine.processNextItem(ctx) {
	}
}

func (engine *Impl) processNextItem(ctx context.Context) bool {
	// Wait until there is a new item in the working queue
	key, quit := engine.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer engine.queue.Done(key)

	// Invoke the method containing the business logic
	err := engine.syncHandler(ctx, key.(string))
	engine.handleErr(err, key)
	return true
}

func (engine *Impl) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		engine.queue.Forget(key)
		return
	}

	// This controller retries maxRetries times if something goes wrong. After that, it stops trying.
	if engine.queue.NumRequeues(key) < engine.maxRetries {
		applog.Infof("Error syncing event %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		engine.queue.AddRateLimited(key)
		return
	}

	engine.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	applog.Infof("Dropping event %q out of the queue: %v", key, err)
}

func (engine *Impl) syncHandler(ctx context.Context, eventId string) error {
	event, err := engine.repository.Get(ctx, eventId)
	if err != nil {
		return err
	}
	return engine.processEvent(ctx, event)
}

// markEventRunning return true if event is marked running
func (engine *Impl) markEventRunning(event *Event) bool {
	engine.Lock()
	defer engine.Unlock()

	if engine.runningSet.Has(event.EventID) {
		return false
	}

	engine.runningSet.Insert(event.EventID)
	return true
}

func (engine *Impl) unmarkEventRunning(event *Event) {
	engine.Lock()
	defer engine.Unlock()

	engine.runningSet.Delete(event.EventID)
}

// Option options of impl
type Option func(impl *Impl)

// WithMaxRetries set max retry
func WithMaxRetries(retry int) Option {
	return func(impl *Impl) {
		impl.maxRetries = retry
	}
}

// WithSyncPeriod set sync period
func WithSyncPeriod(duration time.Duration) Option {
	return func(impl *Impl) {
		impl.syncPeriod = duration
	}
}

// WithBatchSize set batch size
func WithBatchSize(batchSize int) Option {
	return func(impl *Impl) {
		impl.batchSize = batchSize
	}
}
