package main

import (
	"flag"
	"log"
	"net/http"
	"receipt_manager/receipt_upload"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "5000", "Port to run the server on")
	flag.Parse()

	http.HandleFunc("/hello", receiptUploadHandler)
	http.HandleFunc("/test", testReceiptUploadHandler)

	log.Printf("Server running on http://0.0.0.0:%v", port)

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func receiptUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("receipt")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO user and user settings which will determine if to save file ETC
	// for now assuming always save
	receipt_upload.ShouldSaveFile(file, handler)

	receipt_upload.ParseReceipt("./uploads/" + handler.Filename)
}

func testReceiptUploadHandler(w http.ResponseWriter, r *http.Request) {
	receipt_upload.ParseReceipt("./uploads/" + "receipt_test.png")
}
