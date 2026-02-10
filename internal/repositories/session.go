package repositories

import (
	"capecom-pm/internal/domain/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepo struct {
	DB *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{
		DB: db,
	}
}

// Create creates a new session
func (r *SessionRepo) Create(session *models.Session) error {
	if session.UUID == "" {
		session.UUID = uuid.NewString()
	}
	if session.Status == "" {
		session.Status = models.SessionStatusActive
	}

	now := time.Now()
	session.LastUsedAt = now

	return r.DB.Create(session).Error
}

// GetByUUID retrieves a session by UUID
func (r *SessionRepo) GetByUUID(sessionUUID string) (*models.Session, error) {
	var session models.Session
	err := r.DB.Where("uuid = ? AND deleted_at IS NULL", sessionUUID).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &session, err
}

// GetByHashedToken retrieves a session by token (refresh_token_hash)
func (r *SessionRepo) GetByHashedToken(refreshToken string) (*models.Session, error) {
	var session models.Session
	err := r.DB.Where("refresh_token_hash = ? AND deleted_at IS NULL", refreshToken).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &session, err
}

// GetActiveByUserID retrieves all active sessions for a user
func (r *SessionRepo) GetActiveByUserID(userID uint64) ([]models.Session, error) {
	var sessions []models.Session
	err := r.DB.Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, models.SessionStatusActive).
		Order("last_used_at DESC").
		Find(&sessions).Error

	return sessions, err
}

// Update updates an existing session
func (r *SessionRepo) Update(session *models.Session) error {
	session.UpdatedAt = time.Now()
	return r.DB.Save(session).Error
}

// UpdateStatus updates the status of a session
func (r *SessionRepo) UpdateStatus(sessionUUID string, status string) error {
	return r.DB.Model(&models.Session{}).
		Where("uuid = ? AND deleted_at IS NULL", sessionUUID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// UpdateLastUsed updates the last_used_at timestamp
func (r *SessionRepo) UpdateLastUsed(sessionUUID string) error {
	return r.DB.Model(&models.Session{}).
		Where("uuid = ? AND deleted_at IS NULL", sessionUUID).
		Update("last_used_at", time.Now()).Error
}

// RevokeSession revokes a session by UUID
func (r *SessionRepo) RevokeSession(sessionUUID string) error {
	return r.UpdateStatus(sessionUUID, models.SessionStatusRevoked)
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *SessionRepo) RevokeAllUserSessions(userID uint64) error {
	return r.DB.Model(&models.Session{}).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, models.SessionStatusActive).
		Updates(map[string]interface{}{
			"status":     models.SessionStatusRevoked,
			"updated_at": time.Now(),
		}).Error
}

// DeleteExpiredSessions soft deletes expired sessions
func (r *SessionRepo) DeleteExpiredSessions() error {
	now := time.Now()
	return r.DB.Model(&models.Session{}).
		Where("refresh_expires_at < ? AND deleted_at IS NULL", now).
		Update("deleted_at", now).Error
}
