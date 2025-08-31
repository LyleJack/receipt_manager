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

	fs := http.FileServer(http.Dir("./static_page"))
	http.Handle("/", fs)

	http.HandleFunc("/test", testReceiptUploadHandler)
	http.HandleFunc("/upload", receiptUploadHandler)

	log.Printf("Server starting up, go to http://localhost:%v/static_front_end.html to upload your receipt", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func receiptUploadHandler(w http.ResponseWriter, r *http.Request) {
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
	err = receipt_upload.ShouldSaveFile(file, handler)
	if err != nil {
		log.Printf("Error saving the file: %v\n", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	receipt, err := receipt_upload.ParseReceipt("./uploads/" + handler.Filename)
	if err != nil {
		log.Printf("Error parsing the receipt: %v\n", err)
		http.Error(w, "Failed to parse receipt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(receipt)

}

func testReceiptUploadHandler(w http.ResponseWriter, r *http.Request) {
	receipt, err := receipt_upload.ParseReceipt("./uploads/" + "receipt_test.png")
	if err != nil {
		log.Printf("Error parsing the receipt: %v\n", err)
		http.Error(w, "Failed to parse receipt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(receipt)
}
