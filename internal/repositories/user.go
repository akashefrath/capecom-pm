package repositories

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) UserJoinRaw(findBy string, extraQuery string) string {
	userJoin := `
		SELECT 
			u.uuid AS id,
			u.name,
			u.email,
			u.employee_id,
			u.created_at ,
			u.updated_at ,
		    g.uuid AS user_group_id,
			g.name AS user_group_name,
		    de.uuid AS designation_id,
			de.designation AS designation_name,
		    d.uuid AS department_id,
		    d.department AS department_name
		     
		
		FROM users u
		LEFT JOIN user_groups g ON g.id = u.group_id
		LEFT JOIN designations de ON de.id = u.designation_id
		LEFT JOIN  departments d on d.id = u.department_id
		%s
	`
	if findBy != "" {
		userJoin = fmt.Sprintf(userJoin, ` WHERE %s = ? `)
		return fmt.Sprintf(userJoin, findBy) + extraQuery
	} else {
		return fmt.Sprintf(userJoin, "ORDER BY u.name ASC") + extraQuery
	}

}

func (r *UserRepo) UserRolesJoinRaw(findBy string) string {
	roleJoin := `SELECT 
			r.uuid AS id, 
			r.name AS name 
			FROM users u 
			LEFT JOIN  user_roles ur on ur.user_id = u.id
			LEFT JOIN  roles r on r.id = ur.role_id
			WHERE %s = ?
	`
	return fmt.Sprintf(roleJoin, findBy)
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		userIn := tx.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user).Error
		if userIn != nil {
			return userIn
		}
		var count int64 = 0
		err := tx.Raw("SELECT COUNT(*) FROM user_roles WHERE user_id = ? AND (role_id = ? OR role_id = ?) ", user.ID, 1, 2).Count(&count).Error

		if err == nil && count > 0 {
			user.IsAdmin = true
		}
		return err

	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepo) FindByUuid(uuid string) (*models.User, error) {
	var user models.User

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		userIn := tx.Raw("SELECT * FROM users WHERE uuid = ?", uuid).Scan(&user).Error
		if userIn != nil {
			return userIn
		}

		err := tx.Raw("SELECT * FROM user_roles WHERE user_id = ? AND (role_id = ? OR role_id = ?) ", user.ID, 1, 2).Error

		if err == nil {
			user.IsAdmin = true
		}
		return err

	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepo) FindByUUID(uuid string) (*dto.UserResponse, error) {

	var result dto.UserResponse
	var roles []dto.UtilNameId
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Raw(r.UserJoinRaw("u.uuid", ""), uuid).Scan(&result).Error

		err = tx.Raw(r.UserRolesJoinRaw("u.uuid"), uuid).Scan(&roles).Error

		result.Roles = roles
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return err
	})
	if result.Id == "" && err == nil {
		err = domainerrors.ErrUserNotFound
	}
	return &result, err
}
func (r *UserRepo) FindByID(id uint64) (*dto.UserResponse, error) {

	var result dto.UserResponse
	var roles []dto.UtilNameId

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Raw(r.UserJoinRaw("u.id", ""), id).Scan(&result).Error

		err = tx.Raw(r.UserRolesJoinRaw("u.id"), id).Scan(&roles).Error
		result.Roles = roles
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})

	return &result, err
}

func (r *UserRepo) GetUsers(pagination common.Pagination) (*dto.ListWithMeta, error) {

	var result []dto.UserResponse
	result = make([]dto.UserResponse, 0)
	var total int64 = 0

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		err := r.DB.Raw(`SELECT COUNT(*) FROM users`).Scan(&total).Error
		err = r.DB.Raw(r.UserJoinRaw("", pagination.BuildPaginationQuery())).Scan(&result).Error
		return err
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

func (r *UserRepo) FindByIDWithTx(id uint64, tx *gorm.DB) (*dto.UserResponse, error) {

	var result dto.UserResponse
	var roles []dto.UtilNameId

	err := tx.Raw(r.UserJoinRaw("u.id", ""), id).Scan(&result).Error

	err = tx.Raw(r.UserRolesJoinRaw("u.id"), id).Scan(&roles).Error
	result.Roles = roles
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &result, err
}

func (r *UserRepo) GetActiveUserIDByUuid(uuid string) *int64 {
	var id int64
	r.DB.Raw(`SELECT id FROM users WHERE uuid = ?`, uuid).Scan(&id)
	return &id

}

func (r *UserRepo) GetActiveUserUuidByID(id int64) *string {
	var uuid_ string
	r.DB.Raw("SELECT uuid FROM users WHERE id = ?", id).Scan(&uuid_)
	return &uuid_

}

func (r *UserRepo) CreateUserWithRoles(user *models.User, roleIDs []int64) (*dto.UserResponse, error) {
	var results *dto.UserResponse
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Create user
		var mysqlErr *mysql.MySQLError
		if err := tx.Create(user).Error; err != nil {
			if errors.As(err, &mysqlErr) {
				if mysqlErr.Number == 1062 {
					return domainerrors.NewWithCode(http.StatusConflict, domainerrors.ErrDuplicateEmail.Error(), "repo", "CreateUserWithRoles")
				}
			}

			return err
		}
		if len(roleIDs) > 0 {
			userRoles := make([]map[string]any, len(roleIDs))
			for i, roleID := range roleIDs {

				userRoles[i] = map[string]any{
					"uuid":       uuid.NewString(),
					"user_id":    user.ID,
					"role_id":    roleID,
					"status":     "active",
					"created_by": user.CreatedBy,
				}
			}
			if err := tx.Table("user_roles").Create(userRoles).Error; err != nil {
				return err
			}
		}
		result, err := r.FindByIDWithTx(user.ID, tx)
		if err != nil {
			return err
		}
		results = result
		return nil
	})

	return results, err
}

func (r *UserRepo) FindByUuidAndMailIsAdmin(uuid string, roleID int) (bool, error) {
	var userId int64
	var isAdmin = false

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		userIn := tx.Raw("SELECT id FROM users WHERE uuid = ? OR email = ?", uuid, uuid).Scan(&userId).Error
		if userIn != nil {
			return userIn
		}

		var count int64 = 0
		err := tx.Raw("SELECT count(*) FROM user_roles WHERE user_id = ? AND (role_id = ? ) ", userId, roleID).Count(&count).Error

		if err == nil && count > 0 {
			isAdmin = true
		}
		return err

	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return isAdmin, nil
	}

	return isAdmin, err
}

func (r *UserRepo) IsAdmin(uuid string) (bool, error) {
	result, err := r.FindByUuidAndMailIsAdmin(uuid, 1)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (r *UserRepo) IsManager(uuid string) (bool, error) {
	result, err := r.FindByUuidAndMailIsAdmin(uuid, 3)
	if err != nil {
		return false, err
	}

	return result, nil
}
