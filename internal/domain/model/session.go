package models

import "time"

//const (
//	SessionStatusActive   = "active"
//	SessionStatusInactive = "inactive"
//	SessionStatusBlocked  = "blocked"
//	SessionStatusRevoked  = "revoked"
//)

type Session struct {
	BaseModelNoCB

	UserID           int64      `db:"user_id"`
	JTI              string     `db:"jti"`
	RefreshTokenHash string     `db:"refresh_token_hash"`
	RefreshExpiresAt time.Time  `db:"refresh_expires_at"`
	RotatedAt        *time.Time `db:"rotated_at"`
	DeviceID         *string    `db:"device_id"`
	DeviceName       *string    `db:"device_name"`
	UserAgent        *string    `db:"user_agent"`
	IPAddress        *string    `db:"ip_address"`
	LastUsedAt       time.Time  `db:"last_used_at"`
}
