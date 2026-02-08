package mastersreo

import (
	"gorm.io/gorm"
)

type DesignationRepo struct {
	DB *gorm.DB
}

func NewDesignationRepo(db *gorm.DB) *DesignationRepo {
	return &DesignationRepo{
		DB: db,
	}
}

func (r *DesignationRepo) GetIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Table("designations").
		Select("id").
		Where("uuid = ? AND status = ?", uuid, "active").
		Scan(&id).Error
	return id, err
}
