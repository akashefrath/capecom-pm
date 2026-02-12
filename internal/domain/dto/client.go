package dto

import "capecom-pm/internal/domain/models"

type CreateClientRequest struct {
	Name    string  `json:"name" form:"name" binding:"required,min=2,max=150"`
	Email   *string `json:"email" form:"email" binding:"omitempty,email,max=255"`
	Phone   *string `json:"phone" form:"phone" binding:"omitempty,max=20"`
	Address *string `json:"address" form:"address" binding:"omitempty"`
}

type UpdateClientRequest struct {
	Name    *string `json:"name" form:"name" binding:"omitempty,min=2,max=150"`
	Email   *string `json:"email" form:"email" binding:"omitempty,email,max=255"`
	Phone   *string `json:"phone" form:"phone" binding:"omitempty,max=20"`
	Address *string `json:"address" form:"address" binding:"omitempty"`
}

type ClientResponse struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Email   *string `json:"email"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`

	models.BaseResponse
}
