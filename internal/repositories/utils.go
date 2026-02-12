package repositories

import (
	"capecom-pm/internal/domain/dto"

	"gorm.io/gorm"
)

type UtilsRepo struct {
	DB *gorm.DB
}

func NewUtilsRepo(db *gorm.DB) *UtilsRepo {
	return &UtilsRepo{DB: db}
}

func (r *UtilsRepo) GetRoles() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, name FROM roles WHERE status = 'active' AND deleted_at IS NULL ORDER BY name`).Scan(&result).Error
	return result, err
}

func (r *UtilsRepo) GetUserGroups() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, name FROM user_groups WHERE status = 'active' AND deleted_at IS NULL ORDER BY name`).Scan(&result).Error
	return result, err
}

func (r *UtilsRepo) GetDesignations() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, designation AS name FROM designations WHERE status = 'active' AND deleted_at IS NULL ORDER BY designation`).Scan(&result).Error
	return result, err
}

func (r *UtilsRepo) GetDepartments() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, department AS name FROM departments WHERE status = 'active' AND deleted_at IS NULL ORDER BY department`).Scan(&result).Error
	return result, err
}

func (r *UtilsRepo) GetClients() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, name FROM clients WHERE status = 'active' AND deleted_at IS NULL ORDER BY name`).Scan(&result).Error
	return result, err
}

func (r *UtilsRepo) GetTicketTypes() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, name FROM ticket_types WHERE status = 'active' AND deleted_at IS NULL ORDER BY name`).Scan(&result).Error
	return result, err
}

func (r *UtilsRepo) GetUsers() ([]dto.UtilOption, error) {
	var result []dto.UtilOption
	err := r.DB.Raw(`SELECT uuid, name FROM users WHERE status = 'active' AND deleted_at IS NULL ORDER BY name`).Scan(&result).Error
	return result, err
}
