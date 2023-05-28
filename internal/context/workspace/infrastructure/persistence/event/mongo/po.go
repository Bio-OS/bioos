package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event represents an event to be processed
type Event struct {
	EventID string `json:"id" bson:"id"`
	Type    string `json:"type" bson:"type"`
	Status  string `json:"status" bson:"status"`
	// todo record reason why in current status
	Reason      string              `bson:"reason" bson:"reason"`
	Payload     string              `json:"payload" bson:"payload"`
	RetryCount  int                 `json:"retryCount" bson:"retryCount"`
	CreatedAt   primitive.DateTime  `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime  `json:"updatedAt" bson:"updatedAt"`
	ScheduledAt primitive.DateTime  `json:"scheduledAt" bson:"scheduledAt"`
	DeletedAt   *primitive.DateTime `json:"deletedAt" bson:"deletedAt"`
}
