package models

import "github.com/jmoiron/sqlx"

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// GetAll возвращает все услуги без пагинации (для обратной совместимости)
func (r *ServiceRepository) GetAll() ([]Service, error) {
	var services []Service
	query := `SELECT id, name, price, duration FROM services ORDER BY id`
	err := r.db.Select(&services, query)
	return services, err
}

// GetAllWithPagination возвращает услуги с пагинацией
func (r *ServiceRepository) GetAllWithPagination(params PaginationParams) ([]Service, PaginationResult, error) {
	var services []Service

	// Получение общего количества услуг
	total, err := r.Count()
	if err != nil {
		return nil, PaginationResult{}, err
	}

	// Запрос с пагинацией
	query := `SELECT id, name, price, duration
	          FROM services
	          ORDER BY id
	          LIMIT $1 OFFSET $2`
	err = r.db.Select(&services, query, params.Limit, params.GetOffset())
	if err != nil {
		return nil, PaginationResult{}, err
	}

	pagination := NewPaginationResult(total, params)
	return services, pagination, nil
}

// Count возвращает общее количество услуг
func (r *ServiceRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM services`
	err := r.db.Get(&count, query)
	return count, err
}

// Админка
func (r *ServiceRepository) GetById(id int) (*Service, error) {
	var service Service
	query := `SELECT id, name, price, duration FROM services WHERE id = $1`
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
			RETURNING id
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
