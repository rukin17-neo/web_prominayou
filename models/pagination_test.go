package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewPaginationParams тестирует создание параметров пагинации из HTTP запроса
func TestNewPaginationParams(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantPage  int
		wantLimit int
	}{
		{
			name:      "no parameters - defaults",
			url:       "/users",
			wantPage:  DefaultPage,
			wantLimit: DefaultLimit,
		},
		{
			name:      "valid page and limit",
			url:       "/users?page=3&limit=50",
			wantPage:  3,
			wantLimit: 50,
		},
		{
			name:      "only page parameter",
			url:       "/users?page=5",
			wantPage:  5,
			wantLimit: DefaultLimit,
		},
		{
			name:      "only limit parameter",
			url:       "/users?limit=30",
			wantPage:  DefaultPage,
			wantLimit: 30,
		},
		{
			name:      "invalid page - negative",
			url:       "/users?page=-1&limit=20",
			wantPage:  DefaultPage,
			wantLimit: 20,
		},
		{
			name:      "invalid page - zero",
			url:       "/users?page=0&limit=20",
			wantPage:  DefaultPage,
			wantLimit: 20,
		},
		{
			name:      "invalid page - not a number",
			url:       "/users?page=abc&limit=20",
			wantPage:  DefaultPage,
			wantLimit: 20,
		},
		{
			name:      "invalid limit - negative",
			url:       "/users?page=1&limit=-10",
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "invalid limit - zero",
			url:       "/users?page=1&limit=0",
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "invalid limit - not a number",
			url:       "/users?page=1&limit=xyz",
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "limit exceeds maximum - should cap at MaxLimit",
			url:       "/users?page=1&limit=200",
			wantPage:  1,
			wantLimit: MaxLimit,
		},
		{
			name:      "limit at maximum boundary",
			url:       "/users?page=1&limit=100",
			wantPage:  1,
			wantLimit: MaxLimit,
		},
		{
			name:      "very large page number",
			url:       "/users?page=999999",
			wantPage:  999999,
			wantLimit: DefaultLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			params := NewPaginationParams(req)

			if params.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", params.Page, tt.wantPage)
			}
			if params.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", params.Limit, tt.wantLimit)
			}
		})
	}
}

// TestGetOffset тестирует вычисление offset для SQL запросов
func TestGetOffset(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		limit      int
		wantOffset int
	}{
		{
			name:       "first page",
			page:       1,
			limit:      20,
			wantOffset: 0,
		},
		{
			name:       "second page",
			page:       2,
			limit:      20,
			wantOffset: 20,
		},
		{
			name:       "third page",
			page:       3,
			limit:      20,
			wantOffset: 40,
		},
		{
			name:       "page 10 with limit 50",
			page:       10,
			limit:      50,
			wantOffset: 450,
		},
		{
			name:       "page 1 with limit 1",
			page:       1,
			limit:      1,
			wantOffset: 0,
		},
		{
			name:       "page 100 with limit 10",
			page:       100,
			limit:      10,
			wantOffset: 990,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PaginationParams{
				Page:  tt.page,
				Limit: tt.limit,
			}
			got := params.GetOffset()
			if got != tt.wantOffset {
				t.Errorf("GetOffset() = %d, want %d", got, tt.wantOffset)
			}
		})
	}
}

// TestNewPaginationResult тестирует создание результата пагинации
func TestNewPaginationResult(t *testing.T) {
	tests := []struct {
		name            string
		total           int
		page            int
		limit           int
		wantTotalPages  int
		wantHasNext     bool
		wantHasPrevious bool
	}{
		{
			name:            "first page with results",
			total:           100,
			page:            1,
			limit:           20,
			wantTotalPages:  5,
			wantHasNext:     true,
			wantHasPrevious: false,
		},
		{
			name:            "middle page",
			total:           100,
			page:            3,
			limit:           20,
			wantTotalPages:  5,
			wantHasNext:     true,
			wantHasPrevious: true,
		},
		{
			name:            "last page",
			total:           100,
			page:            5,
			limit:           20,
			wantTotalPages:  5,
			wantHasNext:     false,
			wantHasPrevious: true,
		},
		{
			name:            "only one page",
			total:           10,
			page:            1,
			limit:           20,
			wantTotalPages:  1,
			wantHasNext:     false,
			wantHasPrevious: false,
		},
		{
			name:            "empty results",
			total:           0,
			page:            1,
			limit:           20,
			wantTotalPages:  0,
			wantHasNext:     false,
			wantHasPrevious: false,
		},
		{
			name:            "partial last page",
			total:           95,
			page:            5,
			limit:           20,
			wantTotalPages:  5,
			wantHasNext:     false,
			wantHasPrevious: true,
		},
		{
			name:            "exactly divisible",
			total:           100,
			page:            2,
			limit:           50,
			wantTotalPages:  2,
			wantHasNext:     false,
			wantHasPrevious: true,
		},
		{
			name:            "single item per page",
			total:           5,
			page:            3,
			limit:           1,
			wantTotalPages:  5,
			wantHasNext:     true,
			wantHasPrevious: true,
		},
		{
			name:            "large dataset",
			total:           10000,
			page:            50,
			limit:           100,
			wantTotalPages:  100,
			wantHasNext:     true,
			wantHasPrevious: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PaginationParams{
				Page:  tt.page,
				Limit: tt.limit,
			}
			result := NewPaginationResult(tt.total, params)

			if result.Total != tt.total {
				t.Errorf("Total = %d, want %d", result.Total, tt.total)
			}
			if result.Page != tt.page {
				t.Errorf("Page = %d, want %d", result.Page, tt.page)
			}
			if result.Limit != tt.limit {
				t.Errorf("Limit = %d, want %d", result.Limit, tt.limit)
			}
			if result.TotalPages != tt.wantTotalPages {
				t.Errorf("TotalPages = %d, want %d", result.TotalPages, tt.wantTotalPages)
			}
			if result.HasNext != tt.wantHasNext {
				t.Errorf("HasNext = %v, want %v", result.HasNext, tt.wantHasNext)
			}
			if result.HasPrevious != tt.wantHasPrevious {
				t.Errorf("HasPrevious = %v, want %v", result.HasPrevious, tt.wantHasPrevious)
			}
		})
	}
}

// TestPaginationConstants проверяет значения констант
func TestPaginationConstants(t *testing.T) {
	if DefaultPage != 1 {
		t.Errorf("DefaultPage = %d, want 1", DefaultPage)
	}
	if DefaultLimit != 20 {
		t.Errorf("DefaultLimit = %d, want 20", DefaultLimit)
	}
	if MaxLimit != 100 {
		t.Errorf("MaxLimit = %d, want 100", MaxLimit)
	}
}

// TestPaginationIntegration тестирует полный цикл пагинации
func TestPaginationIntegration(t *testing.T) {
	// Создаем запрос с параметрами
	req := httptest.NewRequest(http.MethodGet, "/users?page=3&limit=25", nil)

	// Получаем параметры пагинации
	params := NewPaginationParams(req)

	// Проверяем offset для SQL запроса
	offset := params.GetOffset()
	expectedOffset := (3 - 1) * 25
	if offset != expectedOffset {
		t.Errorf("Offset = %d, want %d", offset, expectedOffset)
	}

	// Создаем результат пагинации
	total := 200
	result := NewPaginationResult(total, params)

	// Проверяем результат
	expectedTotalPages := 8 // 200 / 25 = 8
	if result.TotalPages != expectedTotalPages {
		t.Errorf("TotalPages = %d, want %d", result.TotalPages, expectedTotalPages)
	}

	// На странице 3 из 8 должна быть следующая и предыдущая
	if !result.HasNext {
		t.Error("HasNext should be true for page 3 of 8")
	}
	if !result.HasPrevious {
		t.Error("HasPrevious should be true for page 3")
	}
}

// BenchmarkNewPaginationParams бенчмарк для создания параметров пагинации
func BenchmarkNewPaginationParams(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/users?page=5&limit=50", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPaginationParams(req)
	}
}

// BenchmarkGetOffset бенчмарк для вычисления offset
func BenchmarkGetOffset(b *testing.B) {
	params := PaginationParams{Page: 10, Limit: 20}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		params.GetOffset()
	}
}

// BenchmarkNewPaginationResult бенчмарк для создания результата пагинации
func BenchmarkNewPaginationResult(b *testing.B) {
	params := PaginationParams{Page: 5, Limit: 20}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPaginationResult(1000, params)
	}
}
