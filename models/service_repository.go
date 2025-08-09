package models

import "github.com/jmoiron/sqlx"

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Для клиентов сайта
func (r *ServiceRepository) GetAll() ([]Service, error) {
	var services []Service
	query := `Select id, name, price, duration FROM services ORDER BY id`
	err := r.db.Select(&services, query)
	return services, err
}

// Админка
func (r *ServiceRepository) GetById(id int) (*Service, error) {
	var service Service
	query := `Select id, name, price, duration FROM services WHERE id = $1`
	err := r.db.Get(&service, query, id)
	return &service, err

}

func (r *ServiceRepository) Update(service *Service) error {
	query := `
			UPDATE services
			SET name = $1, price = $2, duration = $3
			WHERE id = $4
	`
	_, err := r.db.Exec(
		query,
		service.Name,
		service.Price,
		service.Duration,
		service.ID,
	)
	return err
}

func (r *ServiceRepository) Create(service *Service) error {
	query := `
			INSERT INTO services (name, price, duration)
			VALUES($1, $2, $3)
			RETURING id
	`
	return r.db.QueryRow(
		query,
		service.Name,
		service.Price,
		service.Duration,
	).Scan(&service.ID)
}

func (r *ServiceRepository) Delete(id int) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
