package services

import (
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	"capecom-pm/internal/storage"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	fileRepo  *repositories.FileRepo
	userRepo  *repositories.UserRepo
	r2Client  *storage.R2Client
	redisRepo *cacherepo.RedisRepo
}

func NewFileService(fileRepo *repositories.FileRepo, userRepo *repositories.UserRepo, r2Client *storage.R2Client, redisRepo *cacherepo.RedisRepo) *FileService {
	return &FileService{
		fileRepo:  fileRepo,
		userRepo:  userRepo,
		r2Client:  r2Client,
		redisRepo: redisRepo,
	}
}

func (s *FileService) CreateFileAndGetUploadURL(userUUID string, req dto.CreateFileRequest) (*dto.CreateFileResponse, error) {

	userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
	if err != nil || userID == nil {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "file_service", "get_user_id")
	}

	fileUUID := uuid.NewString()
	storageKey := fileUUID
	if s.r2Client.FolderName != "" {
		storageKey = fmt.Sprintf("%s/%s.%s", s.r2Client.FolderName, fileUUID, req.MimeType)
	}

	file := &models.File{
		BaseModel:    models.NewBase(nil),
		AuthorID:     *userID,
		StorageKey:   storageKey,
		FileStatus:   models.FileStatusPending,
		OriginalName: &req.FileName,
		MimeType:     &req.MimeType,
		Storage:      models.StorageR2,
	}
	file.UUID = fileUUID

	if err := s.fileRepo.Create(file); err != nil {
		return nil, domainerrors.NewWithCode(http.StatusInternalServerError, domainerrors.ErrInternal.Error(), "file_service", "create_file")
	}

	presignedURL, err := s.r2Client.GeneratePresignedUploadURL(storageKey, req.MimeType, 15*time.Minute)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusInternalServerError, domainerrors.ErrInternal.Error(), "file_service", "presign_url")
	}

	// Cache file UUID -> storage_key mapping
	_ = s.redisRepo.SetString(context.Background(), fmt.Sprintf("file:%s", fileUUID), storageKey, 30*time.Minute)

	return &dto.CreateFileResponse{
		FileID:    fileUUID,
		UploadURL: presignedURL,
	}, nil
}

// ConfirmUpload verifies that files exist in R2 and updates their status to uploaded.
func (s *FileService) ConfirmUpload(fileIDs []string) ([]dto.FileUploadStatus, error) {
	results := make([]dto.FileUploadStatus, 0, len(fileIDs))

	for _, fileID := range fileIDs {
		status := dto.FileUploadStatus{FileID: fileID}

		file, err := s.fileRepo.FindByUUID(fileID)
		if err != nil || file == nil {
			results = append(results, status)
			continue
		}

		exists, err := s.r2Client.CheckFileExists(file.StorageKey)
		if err != nil {
			results = append(results, status)
			continue
		}
		status.Exists = exists

		if exists && file.FileStatus == models.FileStatusPending {
			if err := s.fileRepo.UpdateFileStatus(fileID, models.FileStatusUploaded); err == nil {
				status.Updated = true
			}
		}

		results = append(results, status)
	}

	return results, nil
}
