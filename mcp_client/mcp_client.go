package mcpclient

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

type ToolInfo struct {
	Name        string
	Description string
	InputSchema string
}

type MCPClient struct {
	ctx     context.Context
	session *mcp.ClientSession
}

func NewMCPClient(serverPath string) (*MCPClient, error) {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout.
	transport := &mcp.CommandTransport{Command: exec.Command(serverPath)}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, err
	}

	return &MCPClient{
		ctx:     ctx,
		session: session,
	}, nil
}

func (c *MCPClient) Close() {
	c.session.Close()
}

func (c *MCPClient) CallToolsList() ([]ToolInfo, error) {
	var tools []ToolInfo

	for tool, err := range c.session.Tools(c.ctx, nil) {
		if err != nil {
			return nil, fmt.Errorf("failed to list tools: %w", err)
		}

		// Convert input schema to JSON string
		var schemaStr string
		if tool.InputSchema != nil {
			schemaJSON, err := json.Marshal(tool.InputSchema)
			if err == nil {
				schemaStr = string(schemaJSON)
			}
		}

		tools = append(tools, ToolInfo{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: schemaStr,
		})
	}

	return tools, nil
}

func (c *MCPClient) CallGreeting(name string, secondName string) {
	fmt.Println("\n=== Calling Greeting Tool ===")

	params := &mcp.CallToolParams{
		Name: "greet",
		Arguments: map[string]any{
			"name":        name,
			"second_name": secondName,
		},
	}

	res, err := c.session.CallTool(c.ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}

	for _, content := range res.Content {
		textContent := content.(*mcp.TextContent).Text
		var output GreetingOutput
		if err := json.Unmarshal([]byte(textContent), &output); err != nil {
			log.Fatalf("Failed to parse response: %v", err)
		}
		fmt.Printf("Greeting: %s\n", output.Greeting)
		fmt.Printf("MetaData: %s\n", output.MetaData)
	}
}

func (c *MCPClient) CallCryptoCurrency(coinID string, currency string) {
	fmt.Println("\n=== Calling Crypto Price Tool ===")

	params := &mcp.CallToolParams{
		Name: "get_crypto_price",
		Arguments: map[string]any{
			"coin_id":  coinID,
			"currency": currency,
		},
	}

	res, err := c.session.CallTool(c.ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}

	for _, content := range res.Content {
		textContent := content.(*mcp.TextContent).Text
		var output CryptoPriceOutput
		if err := json.Unmarshal([]byte(textContent), &output); err != nil {
			log.Fatalf("Failed to parse response: %v", err)
		}
		fmt.Printf("Coin: %s\n", output.CoinID)
		fmt.Printf("Currency: %s\n", output.Currency)
		fmt.Printf("Price: %.2f\n", output.Price)
	}
}

// CallTool calls a generic MCP tool with the given name and arguments
func (c *MCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]any) (*mcp.CallToolResult, error) {
	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	}

	result, err := c.session.CallTool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", toolName, err)
	}

	if result.IsError {
		return nil, fmt.Errorf("tool %s returned an error", toolName)
	}

	return result, nil
}