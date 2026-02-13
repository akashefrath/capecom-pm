package models

import "time"

type TimeEntry struct {
	BaseModel

	TicketID    uint64
	ProjectID   uint64
	UserID      uint64
	WorkDate    time.Time
	Hours       float64 `gorm:"type:decimal(5,2)"`
	Description *string
}
