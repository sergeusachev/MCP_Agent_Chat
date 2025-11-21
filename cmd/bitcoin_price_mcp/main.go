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
	testContext := "You are a helpful assistant with access to various tools. " +
		"CRITICAL RULE: You MUST NEVER count, calculate, or process data yourself if a tool exists for that task. " +
		"ALWAYS use count_characters tool for ANY character counting task - never count manually. " +
		"For general knowledge questions (like historical facts), answer directly. " +
		"For ALL other operations (files, prices, counting, calculations), you MUST use the appropriate tool. " +
		"When a task requires multiple steps, call tools sequentially."

	agentInstance.SetContext(testContext)

	testMessage(agentInstance, "In which year did WW2 start?")
	testMessage(agentInstance, "How much is bitcoin price in usd?")

	/*
	Реализовано два MCP сервера
		Первый имеет инструменты:
		1. search_files
		2. read_files
		3. save_to_file

		Второй имеет инструменты:
		1. text_to_unicode
	*/
	testMessage(agentInstance,
		"Найди все markdown файлы, начинающиеся на test, прочитай их," +
		"соедини в один текст и конвертируй все в последовательность unicode кодов." + 
		"Результат сохрани в файл md_unicodes.txt")
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
	serverPaths := map[string]string{
		"main-server":           "/Users/sergeyusachev/Projects/GoProjects/MCP_Example/mcp_server/myserver",
		"unicode-converter-server": "/Users/sergeyusachev/Projects/GoProjects/MCP_Example/mcp_server/unicode_converter",
	}
	mcpClient, err := mcpclient.NewMCPClient(serverPaths)
	if err != nil {
		fmt.Println("MCP client creation error: ", err)
		os.Exit(1)
	}
	return mcpClient
}