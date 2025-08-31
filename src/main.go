package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"receipt_manager/db"
	"receipt_manager/receipt_upload"
	"strconv"
	"strings"
)

type DBConnection struct {
	DB *sql.DB
}

func main() {
	var port string
	flag.StringVar(&port, "port", "5000", "Port to run the server on")
	flag.Parse()

	var dbConn DBConnection
	var err error

	dbConn.DB, err = db.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	fs := http.FileServer(http.Dir("./static_page"))
	http.Handle("/", fs)

	http.HandleFunc("/test", dbConn.testReceiptUploadHandler)
	http.HandleFunc("/receipts/upload", dbConn.receiptUploadHandler)
	http.HandleFunc("/receipts/get", dbConn.receiptsGetHandler)
	http.HandleFunc("/receipts/get/items/", dbConn.receiptsGetItemsHandler)

	log.Printf("Server starting up, go to http://localhost:%v/static_front_end.html to manage your receipts", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func (dbc DBConnection) receiptUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("receipt")
	if err != nil {
		log.Printf("Error retrieving the file: %v\n", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO user and user settings which will determine if to save file ETC
	// for now assuming always save
	filePath, err := receipt_upload.ShouldSaveFile(file, handler)
	if err != nil {
		log.Printf("Error saving the file: %v\n", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	receiptJSON, receipt, err := receipt_upload.ParseReceipt(filePath)
	if err != nil {
		log.Printf("Error parsing the receipt: %v\n", err)
		http.Error(w, "Failed to parse receipt", http.StatusInternalServerError)
		return
	}

	db.SaveReceipt(dbc.DB, receipt)

	w.Header().Set("Content-Type", "application/json")
	w.Write(receiptJSON)

}

func (dbc DBConnection) receiptsGetHandler(w http.ResponseWriter, r *http.Request) {
	receipts, err := db.GetAllItems(dbc.DB)
	if err != nil {
		log.Printf("Error fetching receipts: %v\n", err)
		http.Error(w, "Failed to fetch receipts", http.StatusInternalServerError)
		return
	}

	receiptJSON, err := json.Marshal(receipts)
	if err != nil {
		log.Printf("Error converting receipts to JSON: %v\n", err)
		http.Error(w, "Failed to convert receipts to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(receiptJSON)
}

func (dbc DBConnection) receiptsGetItemsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract receipt ID from URL
	// URL format: /receipts/get/item/{id}
	receiptURL := strings.Split(r.URL.Path, "/")
	receiptID, err := strconv.Atoi(receiptURL[len(receiptURL)-1])
	if err != nil {
		log.Printf("Invalid receipt ID for %s: %v\n", receiptURL, err)
		http.Error(w, "Invalid receipt ID", http.StatusBadRequest)
		return
	}
	receipts, err := db.GetItemsByReceipt(dbc.DB, receiptID)
	if err != nil {
		log.Printf("Error fetching receipt items: %v\n", err)
		http.Error(w, "Failed to fetch receipt items", http.StatusInternalServerError)
		return
	}

	receiptJSON, err := json.Marshal(receipts)
	if err != nil {
		log.Printf("Error converting receipts to JSON: %v\n", err)
		http.Error(w, "Failed to convert receipts to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(receiptJSON)
}

func (dbc DBConnection) testReceiptUploadHandler(w http.ResponseWriter, r *http.Request) {
	receiptJSON, _, err := receipt_upload.ParseReceipt("./uploads/" + "receipt_test.png")
	if err != nil {
		log.Printf("Error parsing the receipt: %v\n", err)
		http.Error(w, "Failed to parse receipt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(receiptJSON)
}
