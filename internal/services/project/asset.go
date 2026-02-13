package projectsvc

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	projectrepo "capecom-pm/internal/repositories/project"
	"capecom-pm/internal/storage"
	"net/http"
)

type AssetService struct {
	assetRepo   *projectrepo.AssetRepo
	projectRepo *projectrepo.ProjectRepo
	fileRepo    *repositories.FileRepo
	userRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
	r2Client    *storage.R2Client
}

func NewAssetService(
	assetRepo *projectrepo.AssetRepo,
	projectRepo *projectrepo.ProjectRepo,
	fileRepo *repositories.FileRepo,
	userRepo *repositories.UserRepo,
	redisRepo *cacherepo.RedisRepo,
	r2Client *storage.R2Client,
) *AssetService {
	return &AssetService{
		assetRepo:   assetRepo,
		projectRepo: projectRepo,
		fileRepo:    fileRepo,
		userRepo:    userRepo,
		redisRepo:   redisRepo,
		r2Client:    r2Client,
	}
}

func (s *AssetService) Create(projectUUID string, req dto.CreateProjectAssetRequest, userUUID string) (*dto.ProjectAssetResponse, error) {
	createdBy, err := s.resolveUserID(userUUID)
	if err != nil {
		return nil, err
	}

	projectID, err := s.projectRepo.GetInternalIDByUUID(projectUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrProjectNotFound.Error(), "project_asset_service", "resolve_project")
	}

	var fileID *uint64
	if req.FileUUID != nil && *req.FileUUID != "" {
		fID, err := s.resolveAndVerifyFile(*req.FileUUID)
		if err != nil {
			return nil, err
		}
		fileID = fID
	}

	isPrivate := false
	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}

	asset := &models.ProjectAsset{
		ProjectID:   uint64(projectID),
		Title:       req.Title,
		AssetType:   req.AssetType,
		Description: req.Description,
		FileID:      fileID,
		Content:     req.Content,
		IsPrivate:   isPrivate,
		BaseModel:   models.NewBase(createdBy),
	}

	return s.assetRepo.Create(asset)
}

func (s *AssetService) GetByProjectUUID(projectUUID string, pg common.Pagination) (*dto.ListWithMeta, error) {
	_, err := s.projectRepo.FindByUUID(projectUUID)
	if err != nil {
		return nil, err
	}
	return s.assetRepo.GetByProjectUUID(projectUUID, pg)
}

func (s *AssetService) GetByUUID(uuid string) (*dto.ProjectAssetResponse, error) {
	return s.assetRepo.FindByUUID(uuid)
}

func (s *AssetService) Update(assetUUID string, req dto.UpdateProjectAssetRequest) (*dto.ProjectAssetResponse, error) {
	_, err := s.assetRepo.FindByUUID(assetUUID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.AssetType != nil {
		updates["asset_type"] = *req.AssetType
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.IsPrivate != nil {
		updates["is_private"] = *req.IsPrivate
	}
	if req.FileUUID != nil {
		if *req.FileUUID == "" {
			updates["file_id"] = nil
		} else {
			fID, err := s.resolveAndVerifyFile(*req.FileUUID)
			if err != nil {
				return nil, err
			}
			updates["file_id"] = *fID
		}
	}

	if len(updates) == 0 {
		return s.assetRepo.FindByUUID(assetUUID)
	}

	return s.assetRepo.Update(assetUUID, updates)
}

func (s *AssetService) Delete(assetUUID string) error {
	return s.assetRepo.Delete(assetUUID)
}

func (s *AssetService) resolveAndVerifyFile(fileUUID string) (*uint64, error) {
	file, err := s.fileRepo.FindByUUID(fileUUID)
	if err != nil || file == nil {
		return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrFileNotFound.Error(), "project_asset_service", "resolve_file")
	}

	if file.FileStatus == models.FileStatusPending {
		exists, err := s.r2Client.CheckFileExists(file.StorageKey)
		if err != nil || !exists {
			return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrFileNotUploaded.Error(), "project_asset_service", "verify_file_r2")
		}
		_ = s.fileRepo.UpdateFileStatus(fileUUID, models.FileStatusUploaded)
	}

	return &file.ID, nil
}

func (s *AssetService) resolveUserID(userUUID string) (*uint64, error) {
	if userUUID == "" {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "project_asset_service", "get_user_id")
	}
	userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
	if err != nil || userID == nil {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "project_asset_service", "get_user_id")
	}
	uid := uint64(*userID)
	return &uid, nil
}
