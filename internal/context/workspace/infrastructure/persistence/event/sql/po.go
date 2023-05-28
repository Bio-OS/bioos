package sql

import (
	"time"

	"gorm.io/gorm"
)

// Event represents an event to be processed
type Event struct {
	EventID string `gorm:"primaryKey;type:varchar(100)"`
	Type    string `gorm:"type:varchar(100)"`
	Status  string `gorm:"type:varchar(100)"`
	// todo record reason why in current status
	Reason      string `gorm:"type:varchar(100)"`
	Payload     string
	RetryCount  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ScheduledAt time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
