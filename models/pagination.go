package models

import (
	"net/http"
	"strconv"
)

// PaginationParams содержит параметры пагинации
type PaginationParams struct {
	Page  int
	Limit int
}

// PaginationResult содержит результаты пагинации
type PaginationResult struct {
	Total       int `json:"total"`
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	TotalPages  int `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

const (
	// DefaultPage значение по умолчанию для страницы
	DefaultPage = 1
	// DefaultLimit количество записей по умолчанию
	DefaultLimit = 20
	// MaxLimit максимальное количество записей на страницу
	MaxLimit = 100
)

// NewPaginationParams создает параметры пагинации из HTTP запроса
func NewPaginationParams(r *http.Request) PaginationParams {
	page := DefaultPage
	limit := DefaultLimit

	// Получение параметра page из query string
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Получение параметра limit из query string
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Ограничение максимального значения limit
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset вычисляет offset для SQL запроса
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// NewPaginationResult создает результат пагинации
func NewPaginationResult(total int, params PaginationParams) PaginationResult {
	totalPages := total / params.Limit
	if total%params.Limit != 0 {
		totalPages++
	}

	return PaginationResult{
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrevious: params.Page > 1,
	}
}
