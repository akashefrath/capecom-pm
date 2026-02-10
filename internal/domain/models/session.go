package models

import "time"

const (
	SessionStatusActive   = "active"
	SessionStatusInactive = "inactive"
	SessionStatusBlocked  = "blocked"
	SessionStatusRevoked  = "revoked"
)

type Session struct {
	BaseModelNoCB
	UserID           int64
	JTI              string
	RefreshTokenHash string
	RefreshExpiresAt time.Time
	RotatedAt        *time.Time
	Status           string
	DeviceID         *string
	DeviceName       *string
	UserAgent        *string
	IPAddress        *string
	LastUsedAt       time.Time
}
