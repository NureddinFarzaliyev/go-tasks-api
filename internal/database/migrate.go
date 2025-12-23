package database

import "database/sql"

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		completed_at DATETIME,
		created_at DATETIME NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}
