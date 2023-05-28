package sql

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type eventRepository struct {
	db             *gorm.DB
	dequeueTimeout time.Duration
	runningTimeout time.Duration
}

var _ eventbus.EventRepository = &eventRepository{}

func NewEventRepository(ctx context.Context, db *gorm.DB, expire, runningTimeout time.Duration) (eventbus.EventRepository, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Event{}); err != nil {
		return nil, err
	}
	return &eventRepository{
		db:             db,
		dequeueTimeout: expire,
		runningTimeout: runningTimeout,
	}, nil
}

func (repo *eventRepository) Get(ctx context.Context, id string) (*eventbus.Event, error) {
	event := &Event{
		EventID: id,
	}
	if ret := repo.db.WithContext(ctx).First(event); ret.Error != nil {
		return nil, ret.Error
	}
	return eventPOToEventDO(event), nil
}
func (repo *eventRepository) Save(ctx context.Context, event *eventbus.Event) error {
	e := eventDOToEventPO(event)
	// ref: https://gorm.io/docs/advanced_query.html#FirstOrCreate
	if ret := repo.db.WithContext(ctx).Where(e).Assign(e).FirstOrCreate(&eventbus.Event{}); ret.Error != nil {
		return ret.Error
	}
	return nil
}

func (repo *eventRepository) ListAndLockUnfinishedEvents(ctx context.Context, limit int, eventTypes []string) ([]*eventbus.Event, error) {
	var events []*Event
	now := time.Now()
	var tx = repo.db.WithContext(ctx)
	if len(eventTypes) > 0 {
		tx = tx.Where("type IN ?", eventTypes)
	}
	ret := tx.
		Where("status = ? AND scheduled_at <= ?", eventbus.EventStatusPending, now).
		Or("status = ? and updated_at <= ?", eventbus.EventStatusDequeue, now.Add(-repo.dequeueTimeout)).
		Or("status = ? and updated_at <= ?", eventbus.EventStatusRunning, now.Add(-repo.runningTimeout)).
		Order("updated_at DESC, scheduled_at ASC").
		Limit(limit).
		Find(&events)
	if ret.Error != nil {
		return nil, ret.Error
	}

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		var rowsAffected int64
		for i := range events {
			var updateMap = map[string]interface{}{}
			switch events[i].Status {
			case eventbus.EventStatusPending:
				updateMap["status"] = eventbus.EventStatusDequeue
				updateMap["updated_at"] = now
				updates := tx.Model(&Event{}).
					Where("event_id = ? and status = ?", events[i].EventID, eventbus.EventStatusPending).
					Updates(updateMap)
				err = updates.Error
				rowsAffected = updates.RowsAffected

			case eventbus.EventStatusDequeue:
				updateMap["updated_at"] = now
				updates := tx.Model(&Event{}).
					Where("event_id = ? and status = ? and updated_at <= ?", events[i].EventID, eventbus.EventStatusDequeue, now.Add(-repo.dequeueTimeout)).
					Updates(updateMap)
				err = updates.Error
				rowsAffected = updates.RowsAffected
			case eventbus.EventStatusRunning:
				updateMap["updated_at"] = now
				updates := tx.Model(&Event{}).
					Where("event_id = ? and status = ? and updated_at <= ?", events[i].EventID, eventbus.EventStatusRunning, now.Add(-repo.runningTimeout)).
					Updates(updateMap)
				err = updates.Error
				rowsAffected = updates.RowsAffected
			default:
				err = fmt.Errorf("unsupport status %s", events[i].Status)
			}
			if err != nil {
				return err
			}
			if rowsAffected != 1 {
				return fmt.Errorf("update failed, expected 1 but got %d", rowsAffected)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	result := make([]*eventbus.Event, 0, len(events))
	for _, e := range events {
		event := e
		result = append(result, eventPOToEventDO(event))
	}
	return result, nil
}

func (repo *eventRepository) UpdateStatus(ctx context.Context, event *eventbus.Event, status string) error {
	e := eventDOToEventPO(event)
	if ret := repo.db.WithContext(ctx).Model(e).First(e); ret.Error != nil {
		return nil
	}
	if event.Status == status {
		return nil
	}
	var updateMap = map[string]interface{}{"status": status, "updated_at": time.Now()}
	if ret := repo.db.WithContext(ctx).Model(e).UpdateColumns(updateMap); ret.Error != nil {
		return ret.Error
	}
	return nil
}

func (repo *eventRepository) UpdateRetryCount(ctx context.Context, event *eventbus.Event, retryCount int) error {
	e := eventDOToEventPO(event)
	var updateMap = map[string]interface{}{"retry_count": retryCount, "updated_at": time.Now()}
	if ret := repo.db.WithContext(ctx).Model(e).UpdateColumns(updateMap); ret.Error != nil {
		return ret.Error
	}
	return nil
}

func (repo *eventRepository) Search(ctx context.Context, Filter *eventbus.Filter) ([]*eventbus.Event, error) {
	var events []Event
	db := repo.db.WithContext(ctx)
	if len(Filter.Type) > 0 {
		db = db.Where("type IN ?", Filter.Type)
	}
	if len(Filter.Status) > 0 {
		db = db.Where("status IN ?", Filter.Status)
	}
	if Filter.Payload != "" {
		db = db.Where("payload LIKE ?", "%"+Filter.Payload+"%")
	}

	if err := db.Find(&events).Error; err != nil {
		applog.Errorw("failed to find events", "err", err)
		return nil, err
	}

	eventPOs := make([]*eventbus.Event, len(events))
	for i := range events {
		eventPOs[i] = eventPOToEventDO(&events[i])
	}

	return eventPOs, nil
}
