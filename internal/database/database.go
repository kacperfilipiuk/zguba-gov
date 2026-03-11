package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS found_items (
			id                 VARCHAR(36) PRIMARY KEY,
			municipality_name  VARCHAR(255) NOT NULL,
			municipality_type  VARCHAR(100) NOT NULL,
			municipality_email VARCHAR(255) NOT NULL,
			item_name          VARCHAR(255) NOT NULL,
			item_category      VARCHAR(100) NOT NULL,
			item_date          VARCHAR(20) NOT NULL,
			item_location      VARCHAR(500) NOT NULL,
			item_status        VARCHAR(50) NOT NULL DEFAULT 'available',
			item_description   TEXT,
			pickup_deadline    INT NOT NULL,
			pickup_location    VARCHAR(500) NOT NULL,
			pickup_hours       VARCHAR(255),
			pickup_contact     VARCHAR(255),
			categories         TEXT,
			created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_found_items_municipality (municipality_name),
			INDEX idx_found_items_category (item_category),
			INDEX idx_found_items_name (item_name),
			INDEX idx_found_items_created (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}
