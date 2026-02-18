package utilsrepository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

type Role struct {
	DB *sqlx.DB
}

func NewRole(db *sqlx.DB) *Role {
	return &Role{DB: db}
}

func (r *Role) GetAllActive() ([]dto.RoleResponse, error) {
	var roles []dto.RoleResponse
	q := `SELECT uuid, name FROM roles WHERE deleted_at IS NULL AND status = ?`
	err := r.DB.Select(&roles, q, models.StatusActive)
	return roles, err
}
