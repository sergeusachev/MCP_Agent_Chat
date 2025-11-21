package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"serge.com/mcp-example/agent"
	"serge.com/mcp-example/api/gigachat"
	"serge.com/mcp-example/mcp_client"
)

func main() {
	mcpClient := getMCPClient()
	gigaChatNetworkService := getNetworkService()
	agentInstance, err := agent.NewAgent(gigaChatNetworkService, mcpClient)
	if err != nil {
		fmt.Println("Agent creation error: ", err)
		os.Exit(1)
	}
	testContext := ""
	agentInstance.SetContext(testContext)

	fmt.Println("=== Bitcoin Price Reminder Service ===")
	fmt.Println("Checking Bitcoin price every 1 minute...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Create a ticker that triggers every 1 minutes
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run the first check immediately
	checkBitcoinPrice(agentInstance)

	// Main loop - runs forever until interrupted
	for {
		select {
		case <-ticker.C:
			// Ticker triggered - check Bitcoin price
			checkBitcoinPrice(agentInstance)
		case <-sigChan:
			// Received interrupt signal - exit gracefully
			fmt.Println("\n\nShutting down Bitcoin Price Reminder Service...")
			return
		}
	}
}

func checkBitcoinPrice(agentInstance *agent.Agent) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] Checking Bitcoin price...\n", timestamp)

	question := "What is the current Bitcoin price in USD?"
	answer, err := agentInstance.SendMessage(question)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get Bitcoin price: %v\n\n", err)
		return
	}

	fmt.Printf("\n💰 Bitcoin Price Update:\n%s\n", answer)
	fmt.Println("─────────────────────────────────────────")
	fmt.Println()
}

func getNetworkService() *gigachat.NetworkService {
	networkService, err := gigachat.GetNetworkService()
	if err != nil {
		fmt.Println("Network service creation error: ", err)
		os.Exit(1)
	}
	return networkService
}

func getMCPClient() *mcpclient.MCPClient {
	mcpClient, err := mcpclient.NewMCPClient("/Users/sergeyusachev/Projects/GoProjects/MCP_Example/mcp_server/myserver")
	if err != nil {
		fmt.Println("MCP client creation error: ", err)
		os.Exit(1)
	}
	return mcpClient
}