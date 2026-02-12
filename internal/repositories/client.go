package repositories

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

type ClientRepo struct {
	DB *gorm.DB
}

func NewClientRepo(db *gorm.DB) *ClientRepo {
	return &ClientRepo{DB: db}
}

func (r *ClientRepo) Create(client *models.Client) (*dto.ClientResponse, error) {
	var mysqlErr *mysql.MySQLError
	if err := r.DB.Create(client).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, domainerrors.NewWithCode(http.StatusConflict, domainerrors.ErrDuplicateClient.Error(), "repo", "CreateClient")
		}
		return nil, err
	}

	return r.findByID(client.ID)
}

func (r *ClientRepo) FindByUUID(uuid string) (*dto.ClientResponse, error) {
	var result dto.ClientResponse
	err := r.DB.Raw(r.selectQuery("c.uuid"), uuid).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrClientNotFound
	}
	return &result, nil
}

func (r *ClientRepo) findByID(id uint64) (*dto.ClientResponse, error) {
	var result dto.ClientResponse
	err := r.DB.Raw(r.selectQuery("c.id"), id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrClientNotFound
	}
	return &result, nil
}

func (r *ClientRepo) GetClients(pagination common.Pagination) (*dto.ListWithMeta, error) {
	var result []dto.ClientResponse
	result = make([]dto.ClientResponse, 0)
	var total int64

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`SELECT COUNT(*) FROM clients WHERE deleted_at IS NULL`).Scan(&total).Error; err != nil {
			return err
		}
		return tx.Raw(r.selectQuery("") + ` ORDER BY c.name ASC` + pagination.BuildPaginationQuery()).Scan(&result).Error
	})

	return &dto.ListWithMeta{
		Data: result,
		Meta: dto.PaginationMeta{
			Limit:   pagination.Limit,
			Page:    pagination.Page,
			Total:   total,
			HasMore: pagination.HasMore(total),
		},
	}, err
}

func (r *ClientRepo) Update(uuid string, updates map[string]interface{}) (*dto.ClientResponse, error) {
	res := r.DB.Table("clients").Where("uuid = ? AND deleted_at IS NULL", uuid).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrClientNotFound
	}
	return r.FindByUUID(uuid)
}

func (r *ClientRepo) GetInternalIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Raw(`SELECT id FROM clients WHERE uuid = ? AND deleted_at IS NULL`, uuid).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, domainerrors.ErrClientNotFound
	}
	return id, nil
}

func (r *ClientRepo) selectQuery(whereCol string) string {
	q := `SELECT c.uuid AS id, c.name, c.email, c.phone, c.address, c.created_at, c.updated_at
		FROM clients c WHERE c.deleted_at IS NULL`
	if whereCol != "" {
		q += ` AND ` + whereCol + ` = ?`
	}
	return q
}
