package db

import (
	"database/sql"
	"log"
	"receipt_manager/receipt_upload"
)

// GetAllItems retrieves all items from all receipts.
func GetAllItems(db *sql.DB) ([]receipt_upload.Receipt, error) {
	rows, err := db.Query(`SELECT id, storename, category, date, total, tip FROM receipts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []receipt_upload.Receipt
	for rows.Next() {
		var receipt receipt_upload.Receipt
		if err := rows.Scan(&receipt.ID, &receipt.StoreName, &receipt.Category, &receipt.Date, &receipt.Total, &receipt.Tip); err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	log.Println("Fetched receipts:", receipts)
	return receipts, rows.Err()
}

// GetItemsByReceipt retrieves all items for a specific receipt.
func GetItemsByReceipt(db *sql.DB, receiptID int) ([]receipt_upload.Item, error) {
	rows, err := db.Query(
		`SELECT name, quantity, total_price FROM receipt_items WHERE receipt_id = ?`, receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []receipt_upload.Item
	for rows.Next() {
		var item receipt_upload.Item
		if err := rows.Scan(&item.Name, &item.Quantity, &item.TotalPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
