package receipt_upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	gosseract "github.com/otiai10/gosseract/v2"
	"google.golang.org/genai"
)

const (
	packageName = "receipt_upload"
)

type Receipt struct {
	StoreName string    `json:"store_name"`
	Location  string    `json:"location,omitempty"`
	Date      time.Time `json:"date"`
	Items     []Item    `json:"items"`
	Total     float64   `json:"total"`
	Tip       float64   `json:"tip,omitempty"`
}

type Item struct {
	Quanity    int     `json:"quantity,omitempty"`
	Name       string  `json:"name"`
	TotalPrice float64 `json:"total_price"`
}

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

	receiptText, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	fmt.Println("image text here:", receiptText)

	ctx := context.Background()
	aiClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	// change this to run as a goroutine at a later point

	receipt, err := extractReceipt(ctx, aiClient, receiptText)
	if err != nil {
		log.Fatalf("Failed to extract receipt: %v", err)
	}

	_ = receipt

	// save receipt json in database here

	return "", nil
}

func extractReceipt(ctx context.Context, client *genai.Client, ocrText string) (Receipt, error) {
	// Define the JSON schema for the desired output.
	// This is where you specify the exact structure you want.

	jsonSchema := genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"store_name": {
				Type:        genai.TypeString,
				Description: "The name of the store.",
			},
			"date": {
				Type:        genai.TypeString,
				Description: "The transaction date turned into a suitable time.Time golang type format",
			},
			"total": {
				Type:        genai.TypeNumber,
				Description: "The total amount of the receipt.",
			},
			"items": {
				Type:        genai.TypeArray,
				Description: "A list of items purchased.",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name": {
							Type:        genai.TypeString,
							Description: "The name of the item.",
						},
						"total_price": {
							Type:        genai.TypeNumber,
							Description: "The price of the items.",
						},
						"quantity": {
							Type:        genai.TypeInteger,
							Description: "The quantity of the item purchased.",
						},
					},
				},
			},
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:  genai.TypeArray,
			Items: &jsonSchema,
		},
	}

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: (fmt.Sprintf(`Convert the following receipt text to a JSON object:
					%s`, ocrText)),
				},
			},
		},
	}, config)
	if err != nil {
		return Receipt{}, fmt.Errorf("generate content error: %w", err)
	}

	textVal := resp.Text()
	log.Printf("Full response text: %s", textVal)

	// Unmarshal the raw JSON into our Go struct.
	var receipt []Receipt
	if err := json.Unmarshal([]byte(textVal), &receipt); err != nil {
		return Receipt{}, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	if len(receipt) == 0 {
		return Receipt{}, fmt.Errorf("no receipt data found in response")
	}
	if len(receipt) > 1 {
		log.Printf("Warning: multiple receipt objects found, using the first one")
	}
	return receipt[0], nil
}
