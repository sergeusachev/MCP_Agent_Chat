package main

import (
	"context"
	"fmt"
	"log"

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

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "get_crypto_price", Description: "Getter for crypto price"}, GetCryptoPrice)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}