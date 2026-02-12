package repositories

import (
	"capecom-pm/internal/domain/models"
	"errors"

	"gorm.io/gorm"
)

type FileRepo struct {
	DB *gorm.DB
}

func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{DB: db}
}

func (r *FileRepo) Create(file *models.File) error {
	return r.DB.Create(file).Error
}

func (r *FileRepo) FindByID(id int64) (*models.File, error) {
	var file models.File
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

func (r *FileRepo) FindByUUID(uuid string) (*models.File, error) {
	var file models.File
	err := r.DB.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

func (r *FileRepo) FindByStorageKey(key string) (*models.File, error) {
	var file models.File
	err := r.DB.Where("storage_key = ? AND deleted_at IS NULL", key).First(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &file, err
}

func (r *FileRepo) UpdateFileStatus(uuid string, status string) error {
	result := r.DB.Model(&models.File{}).
		Where("uuid = ? AND deleted_at IS NULL", uuid).
		Update("file_status", status)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *FileRepo) UpdateFileStatusBulk(uuids []string, status string) error {
	result := r.DB.Model(&models.File{}).
		Where("uuid IN ? AND deleted_at IS NULL", uuids).
		Update("file_status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FileExistsResult maps an identifier to its DB file ID (0 if not found).
type FileExistsResult struct {
	Key    string `json:"key"`
	FileID uint64 `json:"file_id"`
}

// FindExistenceByUUIDs takes a list of UUIDs and returns each with its file ID (0 if not found).
func (r *FileRepo) FindExistenceByUUIDs(uuids []string) ([]FileExistsResult, error) {
	var files []models.File
	if err := r.DB.Select("id, uuid").
		Where("uuid IN ? AND deleted_at IS NULL", uuids).
		Find(&files).Error; err != nil {
		return nil, err
	}

	lookup := make(map[string]uint64, len(files))
	for _, f := range files {
		lookup[f.UUID] = f.ID
	}

	results := make([]FileExistsResult, len(uuids))
	for i, u := range uuids {
		results[i] = FileExistsResult{Key: u, FileID: lookup[u]}
	}
	return results, nil
}

// FindExistenceByStorageKeys takes a list of storage keys and returns each with its file ID (0 if not found).
func (r *FileRepo) FindExistenceByStorageKeys(keys []string) ([]FileExistsResult, error) {
	var files []models.File
	if err := r.DB.Select("id, storage_key").
		Where("storage_key IN ? AND deleted_at IS NULL", keys).
		Find(&files).Error; err != nil {
		return nil, err
	}

	lookup := make(map[string]uint64, len(files))
	for _, f := range files {
		lookup[f.StorageKey] = f.ID
	}

	results := make([]FileExistsResult, len(keys))
	for i, k := range keys {
		results[i] = FileExistsResult{Key: k, FileID: lookup[k]}
	}
	return results, nil
}
