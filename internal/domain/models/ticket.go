package models

import "time"

type Ticket struct {
	BaseModel

	ProjectID              uint64
	Code                   string
	Title                  string
	Description            *string
	Branch                 *string
	TicketTypeID           uint64
	AssignedTo             *uint64
	AssignedBy             *uint64
	StartDate              *time.Time
	InternalStartDate      *time.Time
	EndDate                *time.Time
	InternalEndDate        *time.Time
	EstimatedHours         float64 `gorm:"type:decimal(7,2)"`
	InternalEstimatedHours float64 `gorm:"type:decimal(7,2)"`
	LifecycleStatus        string  `gorm:"default:todo"`
	Priority               string  `gorm:"default:medium"`
	ParentID               *uint64
	DueDate                *time.Time
}
