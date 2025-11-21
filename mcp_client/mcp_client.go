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

type serverConnection struct {
	session *mcp.ClientSession
	tools   map[string]bool // Set of tool names provided by this server
}

type MCPClient struct {
	ctx     context.Context
	servers map[string]*serverConnection // Map server name to connection
}

func NewMCPClient(serverPaths map[string]string) (*MCPClient, error) {
	ctx := context.Background()
	servers := make(map[string]*serverConnection)

	for serverName, serverPath := range serverPaths {
		// Create a new client for each server
		client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

		// Connect to the server over stdin/stdout
		transport := &mcp.CommandTransport{Command: exec.Command(serverPath)}
		session, err := client.Connect(ctx, transport, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to server %s: %w", serverName, err)
		}

		// Get list of tools from this server
		tools := make(map[string]bool)
		for tool, err := range session.Tools(ctx, nil) {
			if err != nil {
				return nil, fmt.Errorf("failed to list tools from server %s: %w", serverName, err)
			}
			tools[tool.Name] = true
		}

		servers[serverName] = &serverConnection{
			session: session,
			tools:   tools,
		}
	}

	return &MCPClient{
		ctx:     ctx,
		servers: servers,
	}, nil
}

func (c *MCPClient) Close() {
	for _, server := range c.servers {
		server.session.Close()
	}
}

// findServerForTool finds which server provides the given tool
func (c *MCPClient) findServerForTool(toolName string) (*serverConnection, error) {
	for _, server := range c.servers {
		if server.tools[toolName] {
			return server, nil
		}
	}
	return nil, fmt.Errorf("tool %s not found in any server", toolName)
}

func (c *MCPClient) CallToolsList() ([]ToolInfo, error) {
	var tools []ToolInfo

	// Collect tools from all servers
	for _, server := range c.servers {
		for tool, err := range server.session.Tools(c.ctx, nil) {
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
	}

	return tools, nil
}

func (c *MCPClient) CallGreeting(name string, secondName string) {
	fmt.Println("\n=== Calling Greeting Tool ===")

	// Find server for greet tool
	server, err := c.findServerForTool("greet")
	if err != nil {
		log.Fatalf("Failed to find server for greet tool: %v", err)
	}

	params := &mcp.CallToolParams{
		Name: "greet",
		Arguments: map[string]any{
			"name":        name,
			"second_name": secondName,
		},
	}

	res, err := server.session.CallTool(c.ctx, params)
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

	// Find server for get_crypto_price tool
	server, err := c.findServerForTool("get_crypto_price")
	if err != nil {
		log.Fatalf("Failed to find server for get_crypto_price tool: %v", err)
	}

	params := &mcp.CallToolParams{
		Name: "get_crypto_price",
		Arguments: map[string]any{
			"coin_id":  coinID,
			"currency": currency,
		},
	}

	res, err := server.session.CallTool(c.ctx, params)
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
// Automatically routes to the correct server based on tool name
func (c *MCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]any) (*mcp.CallToolResult, error) {
	// Find which server provides this tool
	server, err := c.findServerForTool(toolName)
	if err != nil {
		return nil, err
	}

	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	}

	result, err := server.session.CallTool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", toolName, err)
	}

	if result.IsError {
		return nil, fmt.Errorf("tool %s returned an error", toolName)
	}

	return result, nil
}