package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/utils/pointer"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

const EventCollection = "events"

type eventRepository struct {
	collection     *mongo.Collection
	dequeueTimeout time.Duration
	runningTimeout time.Duration
}

var _ eventbus.EventRepository = &eventRepository{}

func NewEventRepository(ctx context.Context, db *mongo.Database, dequeueTimeout time.Duration, runningTimeout time.Duration) (eventbus.EventRepository, error) {
	collection := db.Collection(EventCollection)
	if _, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{"id": 1}, Options: options.Index().SetUnique(true)},
	}); err != nil {
		return nil, err
	}
	return &eventRepository{
		collection:     db.Collection(EventCollection),
		dequeueTimeout: dequeueTimeout,
		runningTimeout: runningTimeout,
	}, nil
}

func (repo *eventRepository) Get(ctx context.Context, id string) (*eventbus.Event, error) {
	filter := bson.M{"id": id}
	var event Event
	err := repo.collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, eventbus.ErrEventNotFound
		}
		return nil, err
	}
	return eventPOToEventDO(&event), nil
}

func (repo *eventRepository) UpdateRetryCount(ctx context.Context, event *eventbus.Event, retryCount int) error {
	filter := bson.M{"id": event.EventID}
	update := bson.M{"$set": bson.M{"retryCount": retryCount, "deletedAt": time.Now()}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) Save(ctx context.Context, event *eventbus.Event) error {
	e := eventDOToEventPO(event)
	filter := bson.M{
		"id":          e.EventID,
		"type":        e.Type,
		"status":      e.Status,
		"reason":      e.Reason,
		"payload":     e.Payload,
		"retryCount":  e.RetryCount,
		"createdAt":   e.CreatedAt,
		"updatedAt":   e.UpdatedAt,
		"scheduledAt": e.ScheduledAt,
		"deletedAt":   e.DeletedAt,
	}
	update := bson.M{"$set": filter}
	_, err := repo.collection.UpdateOne(ctx, filter, update, &options.UpdateOptions{
		Upsert: pointer.Bool(true),
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) ListAndLockUnfinishedEvents(ctx context.Context, limit int, eventTypes []string) ([]*eventbus.Event, error) {
	now := time.Now()
	var filter = bson.M{}
	if len(eventTypes) != 0 {
		filter["type"] = bson.M{"$in": eventTypes}
	}
	filter["$or"] = bson.A{
		bson.M{"scheduledAt": bson.M{"$lte": primitive.NewDateTimeFromTime(now)}, "status": eventbus.EventStatusPending},
		bson.M{"updatedAt": bson.M{"$lte": primitive.NewDateTimeFromTime(now.Add(-repo.dequeueTimeout))}, "status": eventbus.EventStatusDequeue},
		bson.M{"updatedAt": bson.M{"$lte": primitive.NewDateTimeFromTime(now.Add(-repo.runningTimeout))}, "status": eventbus.EventStatusRunning},
	}

	findOptions := options.Find().SetSort(bson.M{"updatedAt": -1}).SetSort(bson.M{"scheduledAt": 1}).SetLimit(int64(limit))
	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var events []*Event
	if err = cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	var session mongo.Session
	session, err = repo.collection.Database().Client().StartSession()
	if err != nil {
		return nil, err
	}
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		for i := range events {
			var result *mongo.UpdateResult
			var err error
			switch events[i].Status {
			case eventbus.EventStatusPending:
				updateFilter := bson.M{"scheduledAt": bson.M{"$lte": primitive.NewDateTimeFromTime(now)}, "status": eventbus.EventStatusPending}
				update := bson.M{"status": eventbus.EventStatusDequeue, "updated_at": primitive.NewDateTimeFromTime(now)}
				result, err = repo.collection.UpdateOne(sc, updateFilter, bson.M{"$set": update})
			case eventbus.EventStatusDequeue:
				updateFilter := bson.M{"updatedAt": bson.M{"$lte": primitive.NewDateTimeFromTime(now.Add(-repo.dequeueTimeout))}, "status": eventbus.EventStatusDequeue}
				update := bson.M{"updatedAt": primitive.NewDateTimeFromTime(now)}
				result, err = repo.collection.UpdateOne(sc, updateFilter, bson.M{"$set": update})
			case eventbus.EventStatusRunning:
				updateFilter := bson.M{"updatedAt": bson.M{"$lte": primitive.NewDateTimeFromTime(now.Add(-repo.runningTimeout))}, "status": eventbus.EventStatusRunning}
				update := bson.M{"updatedAt": primitive.NewDateTimeFromTime(now)}
				result, err = repo.collection.UpdateOne(sc, updateFilter, bson.M{"$set": update})
			default:
				err = fmt.Errorf("unsupport status %s", events[i].Status)
			}
			if err != nil {
				return err
			}
			if result.MatchedCount != 1 || result.ModifiedCount != 1 {
				return fmt.Errorf("update failed, expected 1 but got %d", result.MatchedCount)
			}
		}
		return session.CommitTransaction(sc)
	})
	session.EndSession(ctx)
	if err != nil {
		return nil, err
	}
	var result []*eventbus.Event
	for i := range events {
		result = append(result, eventPOToEventDO(events[i]))
	}
	return result, nil
}

func (repo *eventRepository) UpdateStatus(ctx context.Context, event *eventbus.Event, status string) error {
	filter := bson.M{"id": event.EventID}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *eventRepository) Search(ctx context.Context, filter *eventbus.Filter) ([]*eventbus.Event, error) {
	cursor, err := repo.collection.Find(ctx, getFilter(filter))
	if err != nil {
		return nil, err
	}

	var result []*eventbus.Event
	for cursor.Next(ctx) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		eventDO := eventPOToEventDO(&event)
		result = append(result, eventDO)
	}

	return result, nil
}

func getFilter(filter *eventbus.Filter) bson.M {
	res := bson.M{}
	if filter != nil {
		if len(filter.Payload) > 0 {
			res["payload"] = bson.M{"$regex": filter.Payload, "$options": ""}
		}
		if len(filter.Type) > 0 {
			res["type"] = bson.M{"$in": filter.Type}
		}
		if len(filter.Status) > 0 {
			res["status"] = bson.M{"$in": filter.Status}
		}
	}
	return res
}
