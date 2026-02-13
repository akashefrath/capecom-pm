package models

import "time"

type TicketHistory struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	UUID      string `gorm:"type:varchar(36);uniqueIndex;not null"`
	TicketID  uint64 `gorm:"not null"`
	ChangedBy uint64 `gorm:"not null"`
	FieldName string `gorm:"type:varchar(50);not null"`
	OldValue  *string
	NewValue  *string
	Note      *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
