package repositories

import (
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"errors"

	"github.com/go-sql-driver/mysql"
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

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepo) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) CreateUserWithRoles(user *models.User, roleIDs []int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Create user
		var mysqlErr *mysql.MySQLError
		if err := tx.Create(user).Error; err != nil {
			if errors.As(err, &mysqlErr) {
				if mysqlErr.Number == 1062 {
					return domainerrors.ErrDuplicateEmail
				}
			}

			return err
		}

		// Create user_roles entries
		if len(roleIDs) > 0 {
			userRoles := make([]map[string]any, len(roleIDs))
			for i, roleID := range roleIDs {
				var uuid string
				tx.Raw("SELECT UUID()").Scan(&uuid)

				userRoles[i] = map[string]any{
					"uuid":       uuid,
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

		return nil
	})
}
