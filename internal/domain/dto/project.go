package dto

import (
	dto_utils "capecom-pm/internal/domain/dto/utils"
	"capecom-pm/internal/domain/models"
	"time"
)

type CreateProjectRequest struct {
	ProjectName            string  `json:"project_name" form:"project_name" binding:"required,min=2,max=120"`
	ProjectCode            string  `json:"project_code" form:"project_code" binding:"required,min=2,max=120"`
	ClientUUID             *string `json:"client_id" form:"client_id" binding:"omitempty,uuid"`
	LifecycleStatus        *string `json:"lifecycle_status" form:"lifecycle_status" binding:"omitempty,oneof=todo in_progress in_review done closed on_hold"`
	StartDate              *string `json:"start_date" form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	InternalStartDate      *string `json:"internal_start_date" form:"internal_start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate                *string `json:"end_date" form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	InternalEndDate        *string `json:"internal_end_date" form:"internal_end_date" binding:"omitempty,datetime=2006-01-02"`
	EstimatedHours         float64 `json:"estimated_hours" form:"estimated_hours" binding:"min=0"`
	InternalEstimatedHours float64 `json:"internal_estimated_hours" form:"internal_estimated_hours" binding:"min=0"`
	PrimaryRepoURL         *string `json:"primary_repo_url" form:"primary_repo_url" binding:"omitempty,url,max=500"`
}

type UpdateProjectRequest struct {
	ProjectName            *string  `json:"project_name" form:"project_name" binding:"omitempty,min=2,max=120"`
	ProjectCode            *string  `json:"project_code" form:"project_code" binding:"omitempty,min=2,max=120"`
	ClientUUID             *string  `json:"client_id" form:"client_id" binding:"omitempty,uuid"`
	LifecycleStatus        *string  `json:"lifecycle_status" form:"lifecycle_status" binding:"omitempty,oneof=todo in_progress in_review done closed on_hold"`
	StartDate              *string  `json:"start_date" form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	InternalStartDate      *string  `json:"internal_start_date" form:"internal_start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate                *string  `json:"end_date" form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	InternalEndDate        *string  `json:"internal_end_date" form:"internal_end_date" binding:"omitempty,datetime=2006-01-02"`
	EstimatedHours         *float64 `json:"estimated_hours" form:"estimated_hours" binding:"omitempty,min=0"`
	InternalEstimatedHours *float64 `json:"internal_estimated_hours" form:"internal_estimated_hours" binding:"omitempty,min=0"`
	PrimaryRepoURL         *string  `json:"primary_repo_url" form:"primary_repo_url" binding:"omitempty,url,max=500"`
	Status                 *string  `json:"status" form:"status" binding:"omitempty,oneof=active inactive blocked archived"`
}

type UpdateProjectLifecycleRequest struct {
	LifecycleStatus string `json:"lifecycle_status" form:"lifecycle_status" binding:"required,oneof=todo in_progress in_review done closed on_hold"`
}

type ProjectResponse struct {
	InternalID             int64      `json:"-" gorm:"internal_id"`
	Id                     string     `json:"id"`
	ProjectName            string     `json:"project_name"`
	ProjectCode            string     `json:"project_code"`
	ClientID               *string    `json:"client_id"`
	ClientName             *string    `json:"client_name"`
	LifecycleStatus        string     `json:"lifecycle_status"`
	StartDate              *time.Time `json:"start_date"`
	InternalStartDate      *time.Time `json:"internal_start_date"`
	EndDate                *time.Time `json:"end_date"`
	InternalEndDate        *time.Time `json:"internal_end_date"`
	EstimatedHours         float64    `json:"estimated_hours"`
	InternalEstimatedHours float64    `json:"internal_estimated_hours"`
	PrimaryRepoURL         *string    `json:"primary_repo_url"`
	Status                 string     `json:"status"`
	TicketCount            int        `json:"ticket_count"`
	TaskCount              int        `json:"task_count"`
	TotalBookedHours       int        `json:"total_booked_hours"`

	Assets []AssetsResponse `json:"assets,omitempty" gorm:"-"`

	models.BaseResponse
}
type AssetsResponse struct {
	Id          string                 `json:"id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Files       dto_utils.FileResponse `json:"file,omitempty" gorm:"embedded;embeddedPrefix:file_"`
}
