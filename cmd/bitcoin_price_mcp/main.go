package main

import (
	"fmt"
	"os"

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
	testContext := "You are a helpful assistant with access to various tools." +
		"Use tools only when necessary to complete the user's request. For general knowledge questions," + 
		"answer directly. For tasks involving files, cryptocurrency prices, or other tool-specific operations," +
		"use the appropriate tools. When a task requires multiple steps with tools, call them sequentially."

	agentInstance.SetContext(testContext)

	testMessage(agentInstance, "In which year did WW2 start?")
	testMessage(agentInstance, "How much is bitcoin price in usd?")
	testMessage(agentInstance, 
		"Найди все текстовые файлы, прочитай их и коротко перескажи что там. Результат сохрани в файл summary.txt")
	
}

func testMessage(agent *agent.Agent, message string) {
	fmt.Printf("User Message:\n> %s\n\n", message)
	answer, err := agent.SendMessage(message)
	if err != nil {
		fmt.Println("Error getting answer from GigaChat: ", err)
		os.Exit(1)
	}
	fmt.Printf("Agent Answer:\n> %s\n\n", answer)
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