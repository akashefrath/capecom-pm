package models

type ProjectMember struct {
	BaseModel

	UserID         uint64
	ProjectID      uint64
	AllocatedHours float64 `gorm:"type:decimal(7,2)"`
}
