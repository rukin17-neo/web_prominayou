package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		remember_token VARCHAR(255),
		remember_token_expires_at TIMESTAMP,
		reset_token VARCHAR(255),
		reset_token_expires_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_remember_token ON users(remember_token);
	CREATE INDEX IF NOT EXISTS idx_users_reset_token ON users(reset_token);
	`
	_, err := r.db.Exec(query)
	return err
}

// GetAll возвращает всех пользователей без пагинации (для обратной совместимости)
func (r *UserRepository) GetAll() ([]User, error) {
	var users []User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users ORDER BY id`
	err := r.db.Select(&users, query)
	return users, err
}

// GetAllWithPagination возвращает пользователей с пагинацией
func (r *UserRepository) GetAllWithPagination(params PaginationParams) ([]User, PaginationResult, error) {
	var users []User

	// Получение общего количества пользователей
	total, err := r.Count()
	if err != nil {
		return nil, PaginationResult{}, err
	}

	// Запрос с пагинацией
	query := `SELECT id, username, email, password_hash, created_at, updated_at
	          FROM users
	          ORDER BY id
	          LIMIT $1 OFFSET $2`
	err = r.db.Select(&users, query, params.Limit, params.GetOffset())
	if err != nil {
		return nil, PaginationResult{}, err
	}

	pagination := NewPaginationResult(total, params)
	return users, pagination, nil
}

// Count возвращает общее количество пользователей
func (r *UserRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.Get(&count, query)
	return count, err
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	var user User
	query := `SELECT id, username, email, password_hash, remember_token, remember_token_expires_at,
	          reset_token, reset_token_expires_at, created_at, updated_at
	          FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password_hash, remember_token, remember_token_expires_at,
	          reset_token, reset_token_expires_at, created_at, updated_at
	          FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password_hash, remember_token, remember_token_expires_at,
	          reset_token, reset_token_expires_at, created_at, updated_at
	          FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByRememberToken(token string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password_hash, remember_token, remember_token_expires_at,
	          reset_token, reset_token_expires_at, created_at, updated_at
	          FROM users WHERE remember_token = $1 AND remember_token_expires_at > NOW()`
	err := r.db.Get(&user, query, token)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByResetToken(token string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password_hash, remember_token, remember_token_expires_at,
	          reset_token, reset_token_expires_at, created_at, updated_at
	          FROM users WHERE reset_token = $1 AND reset_token_expires_at > NOW()`
	err := r.db.Get(&user, query, token)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(u *User) error {
	query := `INSERT INTO users (username, email, password_hash)
	          VALUES ($1, $2, $3)
	          RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, u.Username, u.Email, u.PasswordHash).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) Update(u *User) error {
	query := `UPDATE users
	          SET username = $1, email = $2, updated_at = NOW()
	          WHERE id = $3`
	_, err := r.db.Exec(query, u.Username, u.Email, u.ID)
	return err
}

func (r *UserRepository) UpdatePassword(id int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, id)
	return err
}

func (r *UserRepository) SetRememberToken(id int, token string, expiresAt time.Time) error {
	query := `UPDATE users
	          SET remember_token = $1, remember_token_expires_at = $2, updated_at = NOW()
	          WHERE id = $3`
	_, err := r.db.Exec(query, sql.NullString{String: token, Valid: true},
		sql.NullTime{Time: expiresAt, Valid: true}, id)
	return err
}

func (r *UserRepository) ClearRememberToken(id int) error {
	query := `UPDATE users
	          SET remember_token = NULL, remember_token_expires_at = NULL, updated_at = NOW()
	          WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) SetResetToken(id int, token string, expiresAt time.Time) error {
	query := `UPDATE users
	          SET reset_token = $1, reset_token_expires_at = $2, updated_at = NOW()
	          WHERE id = $3`
	_, err := r.db.Exec(query, sql.NullString{String: token, Valid: true},
		sql.NullTime{Time: expiresAt, Valid: true}, id)
	return err
}

func (r *UserRepository) ClearResetToken(id int) error {
	query := `UPDATE users
	          SET reset_token = NULL, reset_token_expires_at = NULL, updated_at = NOW()
	          WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
