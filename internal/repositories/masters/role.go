package mastersreo

import (
	"gorm.io/gorm"
)

type RoleRepo struct {
	DB *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		DB: db,
	}
}

func (r *RoleRepo) GetIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Table("roles").
		Select("id").
		Where("uuid = ? AND status = ?", uuid, "active").
		Scan(&id).Error
	return id, err
}

func (r *RoleRepo) GetIDsByUUIDs(uuids []string) ([]int64, error) {
	var ids []int64
	err := r.DB.Table("roles").
		Select("id").
		Where("uuid IN ? AND status = ?", uuids, "active").
		Scan(&ids).Error
	return ids, err
}
