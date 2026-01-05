package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                     int            `db:"id" json:"id"`
	Username               string         `db:"username" json:"username"`
	Email                  string         `db:"email" json:"email"`
	PasswordHash           string         `db:"password_hash" json:"-"`
	RememberToken          sql.NullString `db:"remember_token" json:"-"`
	RememberTokenExpiresAt sql.NullTime   `db:"remember_token_expires_at" json:"-"`
	ResetToken             sql.NullString `db:"reset_token" json:"-"`
	ResetTokenExpiresAt    sql.NullTime   `db:"reset_token_expires_at" json:"-"`
	CreatedAt              time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time      `db:"updated_at" json:"updated_at"`
}
