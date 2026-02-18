package repository

import (
	"time"

	"github.com/akashefrath/capecom-pm/internal/config"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Session struct {
	DB     *sqlx.DB
	Config *config.Config
}

func NewSession(db *sqlx.DB, config *config.Config) *Session {
	return &Session{
		DB:     db,
		Config: config,
	}
}

func (r *Session) CreateSession(userUuid string, jti string, hashedToken []byte) (error, error) {
	if jti == "" {
		jti = uuid.NewString()
	}
	expiresAt := time.Now().Add(time.Hour * time.Duration(r.Config.JWT.RefreshExpireHours))

	q := `
			INSERT INTO sessions (user_id, uuid,jti, refresh_token_hash, refresh_expires_at)
			SELECT u.id,?,?,?,?
			FROM users u
		  	WHERE u.uuid = ?
		`
	_, err := r.DB.Exec(q,
		uuid.NewString(),
		jti,
		hashedToken,
		expiresAt, // 20 days
		userUuid,
	)

	return nil, err
}

func (r *Session) UpdateSession(id int64, hashedToken []byte, jti string) error {
	timeNow := time.Now()
	expiresAt := timeNow.Add(time.Hour * time.Duration(r.Config.JWT.RefreshExpireHours))
	q := `UPDATE sessions SET jti=?, refresh_token_hash = ?, refresh_expires_at = ?, rotated_at = ?,last_used_at =? WHERE id = ?`
	_, err := r.DB.Exec(q, jti, hashedToken, expiresAt, timeNow, timeNow, id)

	return err

}

func (r *Session) GetSessionJTI(hashedToken []byte) (*string, *int64, *string, *bool, error) {

	var jti string
	var id *int64
	var userUuid string
	var userId int64
	var isAdmin bool

	q := `SELECT s.id,s.jti,u.uuid,u.id as userUUID 
         FROM sessions as s 
         INNER JOIN users AS u ON s.user_id = u.id
         WHERE refresh_token_hash = ? AND s.status = ?`
	err := r.DB.QueryRow(q, hashedToken, models.StatusActive).Scan(&id, &jti, &userUuid, &userId)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	q2 := `
    SELECT EXISTS(
        SELECT 1
        FROM user_roles
        WHERE user_id = ?
          AND role_id IN (?, ?)
    )
`
	err = r.DB.Get(&isAdmin, q2, userId, 1, 2)

	return &jti, id, &userUuid, &isAdmin, nil

}

func (r *Session) GetByJTI(jti string) (*models.Session, error) {
	var session models.Session
	q := `SELECT * FROM sessions WHERE jti = ? AND refresh_expires_at > NOW() LIMIT 1 `
	q = r.DB.Rebind(q)
	err := r.DB.Get(&session, q, jti)
	return &session, err

}

func (r *Session) RevokeTokenByJti(jti string) (int64, error) {
	co := int64(0)
	q := `UPDATE sessions SET status = ? WHERE jti = ? AND refresh_expires_at > NOW()  AND status =? LIMIT 1 `
	res, err := r.DB.Exec(q, models.StatusInactive, jti, models.StatusActive)
	if err != nil {
		return co, err
	}

	// 2. Extract the number of rows affected
	count, err := res.RowsAffected()
	if err != nil {
		return co, err
	}

	// 3. Return the count (will be 1 if updated, 0 if jti not found or already inactive)
	return count, nil

}
