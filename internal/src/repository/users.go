package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/domain/common"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	utilsdto "github.com/akashefrath/capecom-pm/internal/domain/dto/utils"
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	DB   *sqlx.DB
	DBTx *database.Database
}

func NewUser(db *sqlx.DB, dbTx *database.Database) *User {
	return &User{
		DB:   db,
		DBTx: dbTx,
	}
}

func (r *User) GetUserMinimalByEmail(email string) (*models.MinimalUser, error) {
	var user models.MinimalUser
	var isAdmin = false
	q1 := `SELECT id,uuid,status,password_hash,is_admin FROM users WHERE email = ? AND deleted_at IS NULL AND status = ?`
	err := r.DB.Get(&user, q1, email, models.StatusActive)
	q2 := `
    SELECT EXISTS(
        SELECT 1
        FROM user_roles
        WHERE user_id = ?
          AND role_id IN (?, ?)
    )
`
	err = r.DB.Get(&isAdmin, q2, user.ID, 1, 2)
	user.IsAdmin = isAdmin

	if err != nil {
		return nil, domainerrors.ErrUserNotFound
	}

	return &user, err

}

func (r *User) FindUserByID(id int64) (*dto.User, error) {
	var res dto.User

	q1 := `SELECT id, uuid, name, email, is_admin, status 
           FROM users 
           WHERE id = ? LIMIT 1`

	if err := r.DB.Get(&res, q1, id); err != nil {
		return nil, err
	}

	q2 := `SELECT r.uuid, r.name 
           FROM user_roles ur 
           INNER JOIN roles r ON r.id = ur.role_id
           WHERE ur.user_id = ? 
             AND ur.status = ? 
             AND ur.deleted_at IS NULL`

	if err := r.DB.Select(&res.Roles, q2, id, models.StatusActive); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return &res, nil
		}
		return nil, err
	}

	return &res, nil
}

func (r *User) GetAll(pg common.Pagination, filter []common.FilterWithKeys) (*utilsdto.ListWithMeta, error) {
	var users []dto.User
	filterQ, args := common.BuildFilterQuery(filter)
	fArgs := args
	paginationQ, pArgs := pg.BuildPaginationQuery()

	// 1. Fetch all users in one go
	q1 := `SELECT  id, uuid, name, email, is_admin, status 
           FROM users 
           WHERE deleted_at IS NULL` + ` AND ` + filterQ + paginationQ
	newArgs := append(args, pArgs...)
	if err := r.DB.Select(&users, q1, newArgs...); err != nil {
		return nil, err
	}

	//if len(users) == 0 {
	//	return users, nil
	//}

	// 2. Collect all User IDs into a slice
	userIDs := make([]int64, len(users))
	userMap := make(map[int64]*dto.User) // Map for quick lookup
	for i := range users {
		userIDs[i] = int64(users[i].ID)
		userMap[int64(users[i].ID)] = &users[i]
		// Initialize the Roles slice so it's [] instead of null in JSON
		users[i].Roles = []dto.IDName{}
	}

	// 3. Fetch ALL roles for ALL users in ONE query
	q2, args, err := sqlx.In(`
        SELECT ur.user_id, r.uuid, r.name 
        FROM user_roles ur 
        INNER JOIN roles r ON r.id = ur.role_id
        WHERE ur.user_id IN (?) 
          AND ur.status = ? 
          AND ur.deleted_at IS NULL`, userIDs, models.StatusActive)

	if err != nil {
		return nil, err
	}

	// We need an auxiliary struct to capture the user_id from the join
	type roleRow struct {
		UserID int64 `db:"user_id"`
		dto.IDName
	}
	var rows []roleRow

	if err := r.DB.Select(&rows, r.DB.Rebind(q2), args...); err != nil {
		return nil, err
	}

	// 4. Map the roles back to the correct users in memory
	for _, row := range rows {
		if user, exists := userMap[row.UserID]; exists {
			user.Roles = append(user.Roles, row.IDName)
		}
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND ` + filterQ

	err = r.DB.Get(&total, countQuery, fArgs...)
	if err != nil {
		return nil, err
	}

	data := utilsdto.ListWithMeta{
		Meta: &utilsdto.PaginationMeta{
			Limit:   pg.Limit,
			Page:    pg.Page,
			Total:   total,
			HasMore: pg.HasMore(total),
		},
		Data:         users,
		AppliedFiler: common.FilterApplied(filter),
	}

	return &data, nil
}

func (r *User) GetActiveUserUuidByID(id int64) *string {
	var userUuid string = ""
	q := `SELECT uuid FROM users WHERE id = ? AND status`
	_ = r.DB.Get(&userUuid, q, id, models.StatusActive)
	return &userUuid

}

func (r *User) GetActiveUserIDByUuid(uuid string) *int64 {
	var id int64
	q := `SELECT id FROM users WHERE uuid = ? AND status`
	_ = r.DB.Get(&id, q, id, models.StatusActive)
	return &id
}

func (r *User) FindUserStatus(id *int64) (*string, error) {
	var status string
	q := `SELECT status FROM users WHERE id = ?`
	_ = r.DB.Get(&status, q, id)
	return &status, nil

}

func (r *User) GetByID(uuid string) (*dto.User, error) {
	var user dto.User
	q := `SELECT id, uuid, name, email, is_admin, status 
           FROM users 
           WHERE uuid = ? AND deleted_at IS NULL`
	if err := r.DB.Get(&user, q, uuid); err != nil {
		return nil, err
	}

	user.Roles = []dto.IDName{}
	q2 := `SELECT r.uuid, r.name 
           FROM user_roles ur 
           INNER JOIN roles r ON r.id = ur.role_id
           WHERE ur.user_id = ? 
             AND ur.status = ? 
             AND ur.deleted_at IS NULL`
	if err := r.DB.Select(&user.Roles, q2, user.ID, models.StatusActive); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return &user, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *User) Create(createdBy int64, req dto.CreateUserRequest) (*int64, error) {
	var userID int64
	var roleID []int64
	err := r.DBTx.WithTx(context.Background(), func(tx *sqlx.Tx) error {
		// insert user
		q := `INSERT INTO users (uuid, name, email, phone, country_code, password_hash, employee_id, created_by,status) 
          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := tx.Exec(q, uuid.NewString(), req.Name, req.Email, req.Phone, nil, req.Password, req.EmployeeID, createdBy, models.StatusActive)
		if err != nil {
			return err
		}
		userID, err = result.LastInsertId()
		if err != nil {
			return err
		}
		// get user roles by uuid
		query, args, err := sqlx.In(
			`SELECT id FROM roles WHERE uuid IN (?)`,
			req.Roles,
		)
		if err != nil {
			return err
		}

		// IMPORTANT for MySQL
		query = tx.Rebind(query)

		err = tx.Select(&roleID, query, args...)
		if err != nil {
			return err
		}

		if roleID == nil || len(roleID) == 0 || len(roleID) != len(req.Roles) {
			return domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.RoleNotFound.Error(), "create_user", "check_role_id")
		}
		///// load user roles
		for _, role := range roleID {
			q2 := `INSERT INTO user_roles (uuid,user_id, role_id) VALUES (?,?, ?)`
			_, err = tx.Exec(q2, uuid.NewString(), userID, role)
			if err != nil {
				return err
			}
		}

		return err
	})

	return &userID, err
}
