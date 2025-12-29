package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Master struct {
	ID        int       `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	PhotoURL  string    `db:"photo_url" json:"photo_url"`
	PhotoData []byte    `db:"photo_data" json:"-"`          // бинарные данные фотографии
	PhotoType string    `db:"photo_type" json:"photo_type"` // MIME тип фотографии
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type MastersRepository struct {
	db *sqlx.DB
}

func NewMastersRepository(db *sqlx.DB) *MastersRepository {
	return &MastersRepository{db: db}
}

func (r *MastersRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS masters (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		photo_url TEXT,
		photo_data BYTEA,
		photo_type TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *MastersRepository) GetAll() ([]Master, error) {
	var masters []Master
	query := `SELECT id, first_name, last_name, photo_url, photo_data, photo_type, created_at FROM masters ORDER BY id`
	if err := r.db.Select(&masters, query); err != nil {
		return nil, err
	}
	return masters, nil
}

func (r *MastersRepository) GetByID(id int) (*Master, error) {
	var m Master
	query := `SELECT id, first_name, last_name, photo_url, photo_data, photo_type, created_at FROM masters WHERE id = $1`
	if err := r.db.Get(&m, query, id); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MastersRepository) Create(m *Master) error {
	query := `INSERT INTO masters (first_name, last_name, photo_url, photo_data, photo_type) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(query, m.FirstName, m.LastName, m.PhotoURL, m.PhotoData, m.PhotoType).Scan(&m.ID, &m.CreatedAt)
}

func (r *MastersRepository) Update(m *Master) error {
	query := `UPDATE masters SET first_name = $1, last_name = $2, photo_url = $3, photo_data = $4, photo_type = $5 WHERE id = $6`
	_, err := r.db.Exec(query, m.FirstName, m.LastName, m.PhotoURL, m.PhotoData, m.PhotoType, m.ID)
	return err
}

func (r *MastersRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM masters WHERE id = $1`, id)
	return err
}

func (r *MastersRepository) UpdatePhotoURL(id int, photoURL string) error {
	_, err := r.db.Exec(`UPDATE masters SET photo_url = $1 WHERE id = $2`, photoURL, id)
	return err
}
