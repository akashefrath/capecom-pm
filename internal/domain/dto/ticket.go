package dto

import (
	"capecom-pm/internal/domain/models"
	"time"
)

type CreateTicketRequest struct {
	Title                  string  `json:"title" form:"title" binding:"required,min=2,max=120"`
	Description            *string `json:"description" form:"description" binding:"omitempty"`
	Branch                 *string `json:"branch" form:"branch" binding:"omitempty,max=120"`
	TicketTypeUUID         string  `json:"ticket_type_id" form:"ticket_type_id" binding:"required,uuid"`
	AssignedToUUID         *string `json:"assigned_to" form:"assigned_to" binding:"omitempty,uuid"`
	StartDate              *string `json:"start_date" form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	InternalStartDate      *string `json:"internal_start_date" form:"internal_start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate                *string `json:"end_date" form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	InternalEndDate        *string `json:"internal_end_date" form:"internal_end_date" binding:"omitempty,datetime=2006-01-02"`
	EstimatedHours         float64 `json:"estimated_hours" form:"estimated_hours" binding:"min=0"`
	InternalEstimatedHours float64 `json:"internal_estimated_hours" form:"internal_estimated_hours" binding:"min=0"`
	Priority               *string `json:"priority" form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ParentUUID             *string `json:"parent_id" form:"parent_id" binding:"omitempty,uuid"`
	DueDate                *string `json:"due_date" form:"due_date" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateTicketRequest struct {
	Title                  *string  `json:"title" form:"title" binding:"omitempty,min=2,max=120"`
	Description            *string  `json:"description" form:"description" binding:"omitempty"`
	Branch                 *string  `json:"branch" form:"branch" binding:"omitempty,max=120"`
	TicketTypeUUID         *string  `json:"ticket_type_id" form:"ticket_type_id" binding:"omitempty,uuid"`
	StartDate              *string  `json:"start_date" form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	InternalStartDate      *string  `json:"internal_start_date" form:"internal_start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate                *string  `json:"end_date" form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	InternalEndDate        *string  `json:"internal_end_date" form:"internal_end_date" binding:"omitempty,datetime=2006-01-02"`
	EstimatedHours         *float64 `json:"estimated_hours" form:"estimated_hours" binding:"omitempty,min=0"`
	InternalEstimatedHours *float64 `json:"internal_estimated_hours" form:"internal_estimated_hours" binding:"omitempty,min=0"`
	Priority               *string  `json:"priority" form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ParentUUID             *string  `json:"parent_id" form:"parent_id" binding:"omitempty,uuid"`
	DueDate                *string  `json:"due_date" form:"due_date" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateTicketLifecycleRequest struct {
	LifecycleStatus string `json:"lifecycle_status" form:"lifecycle_status" binding:"required,oneof=todo in_progress in_review testing done closed reopened"`
}

type UpdateTicketAssigneeRequest struct {
	AssignedToUUID string `json:"assigned_to" form:"assigned_to" binding:"required,uuid"`
}

type TicketResponse struct {
	InternalProjectID      uint64     `json:"-" gorm:"internal_project_id"`
	Id                     string     `json:"id"`
	ProjectID              string     `json:"project_id"`
	Code                   string     `json:"code"`
	Title                  string     `json:"title"`
	Description            *string    `json:"description"`
	Branch                 *string    `json:"branch"`
	TicketTypeID           string     `json:"ticket_type_id"`
	TicketTypeName         string     `json:"ticket_type_name"`
	AssignedTo             *string    `json:"assigned_to"`
	AssignedToName         *string    `json:"assigned_to_name"`
	AssignedBy             *string    `json:"assigned_by"`
	AssignedByName         *string    `json:"assigned_by_name"`
	StartDate              *time.Time `json:"start_date"`
	InternalStartDate      *time.Time `json:"internal_start_date"`
	EndDate                *time.Time `json:"end_date"`
	InternalEndDate        *time.Time `json:"internal_end_date"`
	EstimatedHours         float64    `json:"estimated_hours"`
	InternalEstimatedHours float64    `json:"internal_estimated_hours"`
	LifecycleStatus        string     `json:"lifecycle_status"`
	Priority               string     `json:"priority"`
	ParentID               *string    `json:"parent_id"`
	DueDate                *time.Time `json:"due_date"`
	Status                 string     `json:"status"`

	models.BaseResponse
}
