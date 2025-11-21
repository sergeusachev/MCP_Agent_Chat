package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"serge.com/mcp-example/api/cryptogecko"
)

type Input struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
	SecondName string `json:"second_name" jsonschema:"the second_name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
	MetaData string `json:"meta_data" jsonschema:"meta data about user"`
}

type CryptoPriceInput struct {
	CoinID   string `json:"coin_id" jsonschema:"The ID of the cryptocurrency (e.g. bitcoin, ethereum, solana)"`
	Currency string `json:"currency" jsonschema:"The currency to get the price in (e.g. usd, eur, gbp)"`
}

type CryptoPriceOutput struct {
	CoinID   string  `json:"coin_id" jsonschema:"The cryptocurrency ID"`
	Currency string  `json:"currency" jsonschema:"The currency code"`
	Price    float64 `json:"price" jsonschema:"The current price"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{
		Greeting: "Hi " + input.Name + " " + input.SecondName + "!",
		MetaData: "oh my God!"}, nil
}

func GetCryptoPrice(ctx context.Context, req *mcp.CallToolRequest, input CryptoPriceInput) (
	*mcp.CallToolResult,
	CryptoPriceOutput,
	error,
) {
	// Create CoinGecko client
	client, err := cryptogecko.NewCoinGeckoClient()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to create CoinGecko client: %v", err)},
			},
			IsError: true,
		}, CryptoPriceOutput{}, nil
	}

	// Get the price
	priceResp, err := client.GetCoinPrice(input.CoinID, input.Currency)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to get price: %v", err)},
			},
			IsError: true,
		}, CryptoPriceOutput{}, nil
	}

	// Return successful response
	return nil, CryptoPriceOutput{
		CoinID:   priceResp.CoinID,
		Currency: priceResp.Currency,
		Price:    priceResp.Price,
	}, nil
}

const rootPath = "/Users/sergeyusachev/Projects/GoProjects/MCP_Example"

type SearchFilesInput struct {
	Pattern string `json:"pattern" jsonschema:"File pattern to search for (e.g. *.txt, *.md)"`
}

type SearchFilesOutput struct {
	Files []string `json:"files" jsonschema:"List of found file paths"`
}

type ReadFilesInput struct {
	Files []string `json:"files" jsonschema:"Array of file paths to read"`
}

type ReadFilesOutput struct {
	Text string `json:"text" jsonschema:"Combined text content from all files"`
}

type SaveToFileInput struct {
	Filename string `json:"filename" jsonschema:"Name of the file to save (will be created in root directory)"`
	Text     string `json:"text" jsonschema:"Text content to save to the file"`
}

type SaveToFileOutput struct {
	FilePath string `json:"file_path" jsonschema:"Full path to the saved file"`
	Success  bool   `json:"success" jsonschema:"Whether the save was successful"`
}

func SearchFiles(ctx context.Context, req *mcp.CallToolRequest, input SearchFilesInput) (
	*mcp.CallToolResult,
	SearchFilesOutput,
	error,
) {
	var foundFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		matched, err := filepath.Match(input.Pattern, filepath.Base(path))
		if err != nil {
			return err
		}
		if matched {
			foundFiles = append(foundFiles, path)
		}
		return nil
	})

	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to search files: %v", err)},
			},
			IsError: true,
		}, SearchFilesOutput{}, nil
	}

	return nil, SearchFilesOutput{Files: foundFiles}, nil
}

func ReadFiles(ctx context.Context, req *mcp.CallToolRequest, input ReadFilesInput) (
	*mcp.CallToolResult,
	ReadFilesOutput,
	error,
) {
	var combinedText strings.Builder

	for _, filePath := range input.Files {
		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to read file %s: %v", filePath, err)},
				},
				IsError: true,
			}, ReadFilesOutput{}, nil
		}

		combinedText.WriteString(fmt.Sprintf("\n=== %s ===\n%s\n", filePath, string(content)))
	}

	return nil, ReadFilesOutput{Text: combinedText.String()}, nil
}

func SaveToFile(ctx context.Context, req *mcp.CallToolRequest, input SaveToFileInput) (
	*mcp.CallToolResult,
	SaveToFileOutput,
	error,
) {
	// Create full file path in root directory
	filePath := filepath.Join(rootPath, input.Filename)

	// Write content to file
	err := os.WriteFile(filePath, []byte(input.Text), 0644)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to save file: %v", err)},
			},
			IsError: true,
		}, SaveToFileOutput{Success: false}, nil
	}

	return nil, SaveToFileOutput{
		FilePath: filePath,
		Success:  true,
	}, nil
}

func main() {
	// Create a server with tools
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "get_crypto_price", Description: "Getter for crypto price"}, GetCryptoPrice)
	mcp.AddTool(server, &mcp.Tool{Name: "search_files", Description: "Searches files in filesystem"}, SearchFiles)
	mcp.AddTool(server, &mcp.Tool{Name: "read_files", Description: "Reads files"}, ReadFiles)
	mcp.AddTool(server, &mcp.Tool{Name: "save_to_file", Description: "Saves text content to a file"}, SaveToFile)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}