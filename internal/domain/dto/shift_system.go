package dto

type CreateShiftSystemRequest struct {
	Name          string `json:"name" form:"name" binding:"required"`
	StartTime     string `json:"start_time" form:"start_time" binding:"required"`
	EndTime       string `json:"end_time" form:"end_time" binding:"required"`
	CheckinEarly  int    `json:"checkin_early" form:"checkin_early" binding:"min=0"`
	CheckinLate   int    `json:"checkin_late" form:"checkin_late" binding:"min=0"`
	CheckoutEarly int    `json:"checkout_early" form:"checkout_early" binding:"min=0"`
	CheckoutLate  int    `json:"checkout_late" form:"checkout_late" binding:"min=0"`
	IsOvernight   bool   `json:"is_overnight" form:"is_overnight"`
	IsDefault     bool   `json:"is_default" form:"is_default"`
}

type UpdateShiftSystemRequest struct {
	Name          string `json:"name" form:"name" binding:"required"`
	StartTime     string `json:"start_time" form:"start_time" binding:"required"`
	EndTime       string `json:"end_time" form:"end_time" binding:"required"`
	CheckinEarly  int    `json:"checkin_early" form:"checkin_early" binding:"min=0"`
	CheckinLate   int    `json:"checkin_late" form:"checkin_late" binding:"min=0"`
	CheckoutEarly int    `json:"checkout_early" form:"checkout_early" binding:"min=0"`
	CheckoutLate  int    `json:"checkout_late" form:"checkout_late" binding:"min=0"`
	IsOvernight   bool   `json:"is_overnight" form:"is_overnight"`
	IsDefault     bool   `json:"is_default" form:"is_default"`
}

type ShiftSystemResponse struct {
	BaseModelTop
	Name          string `json:"name" db:"name"`
	StartTime     string `json:"start_time,omitempty" db:"start_time"`
	EndTime       string `json:"end_time,omitempty" db:"end_time"`
	CheckinEarly  int    `json:"checkin_early,omitempty" db:"checkin_early"`
	CheckinLate   int    `json:"checkin_late,omitempty" db:"checkin_late"`
	CheckoutEarly int    `json:"checkout_early,omitempty" db:"checkout_early"`
	CheckoutLate  int    `json:"checkout_late,omitempty" db:"checkout_late"`
	IsOvernight   bool   `json:"is_overnight,omitempty" db:"is_overnight"`
	IsDefault     bool   `json:"is_default,omitempty" db:"is_default"`
	BaseModelBottom
}
