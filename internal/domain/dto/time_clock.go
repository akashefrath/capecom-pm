package dto

import "time"

type TimeClockRequest struct {
	Source    string   `json:"source" form:"source" binding:"required,oneof=mobile biometric admin"`
	Latitude  *float64 `json:"latitude" form:"latitude"`
	Longitude *float64 `json:"longitude" form:"longitude"`
	DeviceID  *string  `json:"device_id" form:"device_id"`
	Remarks   *string  `json:"remarks" form:"remarks"`
}
type TimeClockTimeOutRequest struct {
	Source    string   `json:"source" form:"source" binding:"required,oneof=mobile biometric admin"`
	SummaryID *string  `json:"summary_id" form:"summary_id" binding:"required,uuid"`
	Time      *string  `json:"time" form:"time" binding:"required"`
	Latitude  *float64 `json:"latitude" form:"latitude"`
	Longitude *float64 `json:"longitude" form:"longitude"`
	DeviceID  *string  `json:"device_id" form:"device_id"`
	Remarks   *string  `json:"remarks" form:"remarks"`
}

type TimeClockResponse struct {
	BaseModelTop
	UserID              int64     `json:"-" db:"user_id"`
	AttendanceSummaryID int64     `json:"-" db:"attendance_summary_id"`
	UserUUID            string    `json:"-" db:"user_uuid"`
	LogTime             time.Time `json:"log_time" db:"log_time"`
	LogType             string    `json:"log_type" db:"log_type"`
	Source              string    `json:"source" db:"source"`
	Latitude            *float64  `json:"latitude" db:"latitude"`
	Longitude           *float64  `json:"longitude" db:"longitude"`
	DeviceID            *string   `json:"device_id" db:"device_id"`
	Remarks             *string   `json:"remarks" db:"remarks"`
	BaseModelBottom
}

type AttendanceDetails struct {
	BaseModelTop
	LogTime             time.Time `json:"log_time" db:"log_time"`
	LogType             string    `json:"log_type" db:"log_type"`
	Source              string    `json:"source" db:"source"`
	Remarks             *string   `json:"remarks" db:"remarks"`
	AttendanceSummaryID int64     `json:"-" db:"attendance_summary_id"`

	BaseModelBottom
}

type AttendanceDetailsListWithSummary struct {
	TimeLogs          []AttendanceDetails        `json:"time_logs"`
	AttendanceSummary *AttendanceSummaryResponse `json:"attendance_summary"`
	PendingSummary    *AttendanceSummaryResponse `json:"pending_summary"`
}
