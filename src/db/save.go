package db

import (
	"database/sql"
	"fmt"
	"receipt_manager/receipt_upload"

	_ "github.com/mattn/go-sqlite3"
)

func SaveReceipt(db *sql.DB, r receipt_upload.Receipt) error {
	query := `
		INSERT INTO receipts (storename, date, total, category)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(query, r.StoreName, r.Date, r.Total, r.Category)
	if err != nil {
		return fmt.Errorf("failed to insert receipt: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	for _, item := range r.Items {
		itemQuery := `
			INSERT INTO receipt_items (total_price, name, quantity, receipt_id)
			VALUES (?, ?, ?, ?)
		`
		_, err := db.Exec(itemQuery, item.TotalPrice, item.Name, item.Quantity, id)
		if err != nil {
			return fmt.Errorf("failed to insert receipt item: %w", err)
		}
	}
	return nil
}
