package config

import (
	"database/sql"
	"log"
)

// OptimizeDatabase добавляет индексы и оптимизации для PostgreSQL
func OptimizeDatabase(db *sql.DB) error {
	log.Println("Applying database optimizations...")

	optimizations := []string{
		// Masters table optimization
		`CREATE INDEX IF NOT EXISTS idx_masters_created_at ON masters(created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_masters_name ON masters(first_name, last_name);`,

		// Services table optimization
		`CREATE INDEX IF NOT EXISTS idx_services_id ON services(id);`,
		`CREATE INDEX IF NOT EXISTS idx_services_name ON services(name);`,

		// Existing tables already have these indexes:
		// - users: idx_users_username, idx_users_email, idx_users_remember_token, idx_users_reset_token
		// - sessions: idx_sessions_user_id, idx_sessions_expires_at

		// Analyze tables for query planner optimization
		`ANALYZE masters;`,
		`ANALYZE services;`,
		`ANALYZE users;`,
		`ANALYZE sessions;`,

		// Set connection pool parameters (already configured but documented here)
		// These are set in InitDB() via db.SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime
	}

	for _, query := range optimizations {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: optimization query failed: %v", err)
			// Continue with other optimizations even if one fails
		}
	}

	log.Println("Database optimizations applied successfully")
	return nil
}

// GetDatabaseStats возвращает статистику использования БД
func GetDatabaseStats(db *sql.DB) map[string]interface{} {
	stats := db.Stats()

	return map[string]interface{}{
		"open_connections":  stats.OpenConnections,
		"in_use":            stats.InUse,
		"idle":              stats.Idle,
		"wait_count":        stats.WaitCount,
		"wait_duration_ms":  stats.WaitDuration.Milliseconds(),
		"max_idle_closed":   stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}
