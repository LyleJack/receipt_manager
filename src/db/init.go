package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "receipts.db"

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	createMainTableSQL := `
	CREATE TABLE IF NOT EXISTS receipts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		storename TEXT,
		date DATETIME NOT NULL,
		total REAL NOT NULL,
		category TEXT,
		tip REAL NOT NULL DEFAULT 0
	);`

	_, err = db.Exec(createMainTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create main table: %w", err)
	}

	createItemsTableSQL := `
	CREATE TABLE IF NOT EXISTS receipt_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_price REAL NOT NULL,
		name TEXT
		quantity INTEGER
		receipt_id INTEGER REFERENCES receipts(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(createItemsTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create items table: %w", err)
	}

	return db, nil
}
