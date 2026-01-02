package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(255) PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		expires_at TIMESTAMP NOT NULL,
		ip_address VARCHAR(45),
		user_agent TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SessionRepository) Create(s *Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at, ip_address, user_agent)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING created_at`
	return r.db.QueryRow(query, s.ID, s.UserID, s.ExpiresAt, s.IPAddress, s.UserAgent).
		Scan(&s.CreatedAt)
}

func (r *SessionRepository) GetByID(id string) (*Session, error) {
	var session Session
	query := `SELECT id, user_id, created_at, expires_at, ip_address, user_agent
	          FROM sessions WHERE id = $1 AND expires_at > NOW()`
	err := r.db.Get(&session, query, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetByUserID(userID int) ([]Session, error) {
	var sessions []Session
	query := `SELECT id, user_id, created_at, expires_at, ip_address, user_agent
	          FROM sessions WHERE user_id = $1 AND expires_at > NOW()
	          ORDER BY created_at DESC`
	err := r.db.Select(&sessions, query, userID)
	return sessions, err
}

func (r *SessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SessionRepository) DeleteByUserID(userID int) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *SessionRepository) DeleteExpired() error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}

func (r *SessionRepository) Extend(id string, expiresAt time.Time) error {
	query := `UPDATE sessions SET expires_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, expiresAt, id)
	return err
}
