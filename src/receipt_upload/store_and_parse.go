package receipt_upload

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	gosseract "github.com/otiai10/gosseract/v2"
)

const (
	packageName = "receipt_upload"
)

func ShouldSaveFile(file multipart.File, handler *multipart.FileHeader) error {
	const funcName = "ShouldSaveFile"
	// Create a new file on disk
	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		return fmt.Errorf("[%s].[%s] error: could not create file", packageName, funcName)
	}
	defer dst.Close()

	// Copy uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return fmt.Errorf("[%s].[%s] error: could not save file", packageName, funcName)
	}

	log.Printf("File %s uploaded successfully\n", handler.Filename)

	return nil
}

func ParseReceipt(filePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(filePath)

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	fmt.Println("image text here:", text)

	return text, nil
}
