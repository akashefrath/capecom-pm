package models

import "time"

type Project struct {
	BaseModel

	ProjectName            string
	ProjectCode            string
	ClientID               *uint64
	ClientNameSnapshot     *string
	LifecycleStatus        string `gorm:"default:todo"`
	StartDate              *time.Time
	InternalStartDate      *time.Time
	EndDate                *time.Time
	InternalEndDate        *time.Time
	EstimatedHours         float64 `gorm:"type:decimal(7,2)"`
	InternalEstimatedHours float64 `gorm:"type:decimal(7,2)"`
	PrimaryRepoURL         *string
}
