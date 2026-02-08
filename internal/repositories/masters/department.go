package mastersreo

import (
	"gorm.io/gorm"
)

type DepartmentRepo struct {
	DB *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{
		DB: db,
	}
}

func (r *DepartmentRepo) GetIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Table("departments").
		Select("id").
		Where("uuid = ? AND status = ?", uuid, "active").
		Scan(&id).Error
	return id, err
}
