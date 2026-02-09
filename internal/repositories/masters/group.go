package mastersreo

import (
	"gorm.io/gorm"
)

type GroupRepo struct {
	DB *gorm.DB
}

func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{
		DB: db,
	}
}

func (r *GroupRepo) GetIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Table("user_groups").
		Select("id").
		Where("uuid = ? AND status = ?", uuid, "active").
		Scan(&id).Error
	return id, err
}
