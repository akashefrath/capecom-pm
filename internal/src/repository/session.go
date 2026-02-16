package repository

import (
	"fmt"
	"time"

	"github.com/akashefrath/capecom-pm/internal/config"
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

func (r *Session) CreateSession(userUuid string, jti string, hashedToken string) (error, error) {
	if jti == "" {
		jti = uuid.NewString()
	}
	expiresAt := time.Now().Add(time.Hour * time.Duration(r.Config.JWT.RefreshExpireHours))

	userIDFromUUID := GetUserIDFromUuidQuery()
	_, err := r.DB.Exec(fmt.Sprintf(`
			INSERT INTO sessions (user_id, uuid,jti, refresh_token_hash, refresh_expires_at)
			VALUES ((%s),?,?,?,?)
		`, userIDFromUUID),
		userUuid,
		uuid.NewString(),
		jti,
		hashedToken,
		expiresAt, // 20 days
	)

	return nil, err
}
