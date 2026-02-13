package models

type ProjectAsset struct {
	BaseModel

	ProjectID   uint64
	Title       string
	AssetType   string
	Description *string
	FileID      *uint64
	Content     *string
	IsPrivate   bool `gorm:"default:false"`
}
