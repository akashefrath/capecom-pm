package dto

type TimeClockRequest struct {
	Source    string   `json:"source" form:"source" binding:"required,oneof=mobile biometric admin"`
	Latitude  *float64 `json:"latitude" form:"latitude"`
	Longitude *float64 `json:"longitude" form:"longitude"`
	DeviceID  *string  `json:"device_id" form:"device_id"`
	Remarks   *string  `json:"remarks" form:"remarks"`
}

type TimeClockResponse struct {
	ID        int64    `json:"-" db:"id"`
	UUID      string   `json:"id" db:"uuid"`
	UserID    int64    `json:"-" db:"employee_id"`
	UserUUID  string   `json:"employee_id" db:"user_uuid"`
	LogTime   string   `json:"log_time" db:"log_time"`
	LogType   string   `json:"log_type" db:"log_type"`
	Source    string   `json:"source" db:"source"`
	Latitude  *float64 `json:"latitude" db:"latitude"`
	Longitude *float64 `json:"longitude" db:"longitude"`
	DeviceID  *string  `json:"device_id" db:"device_id"`
	Remarks   *string  `json:"remarks" db:"remarks"`
}
