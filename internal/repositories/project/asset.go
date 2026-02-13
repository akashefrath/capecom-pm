package projectrepo

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type AssetRepo struct {
	DB *gorm.DB
}

func NewAssetRepo(db *gorm.DB) *AssetRepo {
	return &AssetRepo{DB: db}
}

func (r *AssetRepo) Create(asset *models.ProjectAsset) (*dto.ProjectAssetResponse, error) {
	var mysqlErr *mysql.MySQLError
	if err := r.DB.Create(asset).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, domainerrors.NewWithCode(http.StatusConflict, "duplicate_project_asset", "repo", "CreateProjectAsset")
		}
		return nil, err
	}
	return r.findByID(asset.ID)
}

func (r *AssetRepo) FindByUUID(uuid string) (*dto.ProjectAssetResponse, error) {
	var result dto.ProjectAssetResponse
	err := r.DB.Raw(r.selectQuery("pa.uuid"), uuid).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrProjectAssetNotFound
	}
	return &result, nil
}

func (r *AssetRepo) findByID(id uint64) (*dto.ProjectAssetResponse, error) {
	var result dto.ProjectAssetResponse
	err := r.DB.Raw(r.selectQuery("pa.id"), id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrProjectAssetNotFound
	}
	return &result, nil
}

func (r *AssetRepo) GetByProjectUUID(projectUUID string, pg common.Pagination) (*dto.ListWithMeta, error) {
	var result []dto.ProjectAssetResponse
	result = make([]dto.ProjectAssetResponse, 0)
	var total int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`SELECT COUNT(*) FROM project_assets pa
			INNER JOIN projects pr ON pr.id = pa.project_id AND pr.deleted_at IS NULL
			WHERE pa.deleted_at IS NULL AND pr.uuid = ?`, projectUUID).Scan(&total).Error; err != nil {
			return err
		}
		return tx.Raw(r.selectQuery("pr.uuid")+` ORDER BY pa.created_at DESC`+pg.BuildPaginationQuery(), projectUUID).Scan(&result).Error
	})

	return &dto.ListWithMeta{
		Data: result,
		Meta: dto.PaginationMeta{
			Limit:   pg.Limit,
			Page:    pg.Page,
			Total:   total,
			HasMore: pg.HasMore(total),
		},
	}, err
}

func (r *AssetRepo) Update(uuid string, updates map[string]any) (*dto.ProjectAssetResponse, error) {
	res := r.DB.Table("project_assets").Where("uuid = ? AND deleted_at IS NULL", uuid).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrProjectAssetNotFound
	}
	return r.FindByUUID(uuid)
}

func (r *AssetRepo) Delete(uuid string) error {
	res := r.DB.Exec(`UPDATE project_assets SET deleted_at = NOW() WHERE uuid = ? AND deleted_at IS NULL`, uuid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrProjectAssetNotFound
	}
	return nil
}

func (r *AssetRepo) selectQuery(whereCol string) string {
	q := `SELECT pa.uuid AS id, pr.uuid AS project_id, pa.title, pa.asset_type,
		pa.description, f.uuid AS file_id, pa.content, pa.is_private,
		pa.status, pa.created_at, pa.updated_at
		FROM project_assets pa
		INNER JOIN projects pr ON pr.id = pa.project_id AND pr.deleted_at IS NULL
		LEFT JOIN files f ON f.id = pa.file_id AND f.deleted_at IS NULL
		WHERE pa.deleted_at IS NULL`
	if whereCol != "" {
		q += ` AND ` + whereCol + ` = ?`
	}
	return q
}
