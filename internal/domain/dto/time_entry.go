package dto

import (
	"capecom-pm/internal/domain/models"
	"time"
)

type CreateTimeEntryRequest struct {
	WorkDate    string  `json:"work_date" form:"work_date" binding:"required,datetime=2006-01-02"`
	Hours       float64 `json:"hours" form:"hours" binding:"required,min=0.01,max=24"`
	Description *string `json:"description" form:"description" binding:"omitempty"`
}

type UpdateTimeEntryRequest struct {
	WorkDate    *string  `json:"work_date" form:"work_date" binding:"omitempty,datetime=2006-01-02"`
	Hours       *float64 `json:"hours" form:"hours" binding:"omitempty,min=0.01,max=24"`
	Description *string  `json:"description" form:"description" binding:"omitempty"`
}

type TimeEntryResponse struct {
	Id          string    `json:"id"`
	TicketID    string    `json:"ticket_id"`
	TicketCode  string    `json:"ticket_code"`
	ProjectID   string    `json:"project_id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	WorkDate    time.Time `json:"work_date"`
	Hours       float64   `json:"hours"`
	Description *string   `json:"description"`
	Status      string    `json:"status"`

	models.BaseResponse
}

type TicketHistoryResponse struct {
	Id            string    `json:"id"`
	TicketID      string    `json:"ticket_id"`
	ChangedBy     string    `json:"changed_by"`
	ChangedByName string    `json:"changed_by_name"`
	FieldName     string    `json:"field_name"`
	OldValue      *string   `json:"old_value"`
	NewValue      *string   `json:"new_value"`
	Note          *string   `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}
