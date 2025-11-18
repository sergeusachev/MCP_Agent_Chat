package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetingOutput struct {
	Greeting string `json:"greeting"`
	MetaData string `json:"meta_data"`
}

type CryptoPriceOutput struct {
	CoinID   string  `json:"coin_id"`
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}

func main() {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout.
	transport := &mcp.CommandTransport{Command: exec.Command("../mcp_server/myserver")}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Call greeting tool
	callGreeting(ctx, session, "Sergei", "Alexandrovich")

	// Call crypto price tool
	callCryptoCurrency(ctx, session, "bitcoin", "usd")
}

func callGreeting(ctx context.Context, session *mcp.ClientSession, name string, secondName string) {
	fmt.Println("\n=== Calling Greeting Tool ===")

	params := &mcp.CallToolParams{
		Name: "greet",
		Arguments: map[string]any{
			"name":        name,
			"second_name": secondName,
		},
	}

	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}

	for _, c := range res.Content {
		textContent := c.(*mcp.TextContent).Text
		var output GreetingOutput
		if err := json.Unmarshal([]byte(textContent), &output); err != nil {
			log.Fatalf("Failed to parse response: %v", err)
		}
		fmt.Printf("Greeting: %s\n", output.Greeting)
		fmt.Printf("MetaData: %s\n", output.MetaData)
	}
}

func callCryptoCurrency(ctx context.Context, session *mcp.ClientSession, coinID string, currency string) {
	fmt.Println("\n=== Calling Crypto Price Tool ===")

	params := &mcp.CallToolParams{
		Name: "get_crypto_price",
		Arguments: map[string]any{
			"coin_id":  coinID,
			"currency": currency,
		},
	}

	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}

	for _, c := range res.Content {
		textContent := c.(*mcp.TextContent).Text
		var output CryptoPriceOutput
		if err := json.Unmarshal([]byte(textContent), &output); err != nil {
			log.Fatalf("Failed to parse response: %v", err)
		}
		fmt.Printf("Coin: %s\n", output.CoinID)
		fmt.Printf("Currency: %s\n", output.Currency)
		fmt.Printf("Price: %.2f\n", output.Price)
	}
}