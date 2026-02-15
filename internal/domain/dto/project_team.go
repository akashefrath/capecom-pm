package dto

import "capecom-pm/internal/domain/models"

// --- Managers ---

type AssignManagersRequest struct {
	UserUUIDs []string `json:"user_ids"  form:"user_ids" binding:"required,min=1,dive,uuid"`
}

type ProjectManagerResponse struct {
	Id       string `json:"id"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Status   string `json:"status"`

	models.BaseResponse
}

// --- Members ---

type MemberEntry struct {
	UserUUID       string  `json:"user_id" binding:"required,uuid"`
	AllocatedHours float64 `json:"allocated_hours" binding:"min=0"`
}

type AssignMembersRequest struct {
	Members []MemberEntry `json:"members" form:"members" binding:"required,min=1,dive"`
}

type RemoveMembersRequest struct {
	UserUUIDs []string `json:"user_ids" form:"user_ids" binding:"required,min=1,dive,uuid"`
}

type ProjectMemberResponse struct {
	Id             string  `json:"id"`
	UserID         string  `json:"user_id"`
	UserName       string  `json:"user_name"`
	AllocatedHours float64 `json:"allocated_hours"`
	Status         string  `json:"status"`

	models.BaseResponse
}
