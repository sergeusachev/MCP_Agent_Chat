package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"serge.com/mcp-example/api/gigachat"
	"serge.com/mcp-example/mcp_client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Agent struct {
	Model                  string
	temperature            float64
	messages               []gigachat.Message
	gigaChatNetworkService *gigachat.NetworkService
	mcpClient              *mcpclient.MCPClient
	functions              []gigachat.Function
}

func NewAgent(gigaChatNetworkService *gigachat.NetworkService, mcpClient *mcpclient.MCPClient) (*Agent, error) {
	agent := &Agent{
		Model:                  "GigaChat-2",
		temperature:            0.0,
		messages:               []gigachat.Message{},
		gigaChatNetworkService: gigaChatNetworkService,
		mcpClient:              mcpClient,
		functions:              []gigachat.Function{},
	}

	// Get tools from MCP and convert to GigaChat functions
	if err := agent.loadMCPTools(); err != nil {
		return nil, fmt.Errorf("failed to load MCP tools: %w", err)
	}

	return agent, nil
}

func (a *Agent) loadMCPTools() error {
	tools, err := a.mcpClient.CallToolsList()
	if err != nil {
		return err
	}

	for _, tool := range tools {
		// Parse InputSchema from JSON string to object
		var schema any
		if tool.InputSchema != "" {
			if err := json.Unmarshal([]byte(tool.InputSchema), &schema); err != nil {
				return fmt.Errorf("failed to parse schema for tool %s: %w", tool.Name, err)
			}
		}

		a.functions = append(a.functions, gigachat.Function{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  schema,
		})
	}

	return nil
}

func (a *Agent) SetContext(agentContext string) {
	agentContextMessage := gigachat.Message{
		Role:    "system",
		Content: agentContext,
	}

	a.messages = append(a.messages, agentContextMessage)
}

func (a *Agent) SendMessage(message string) (string, error) {
	a.messages = append(a.messages, gigachat.Message{
		Role:    "user",
		Content: message,
	})

	// Call GigaChat with functions
	result, err := a.gigaChatNetworkService.GetCompletion(a.messages, a.Model, a.temperature, a.functions)
	if err != nil {
		return "", fmt.Errorf("failed to get answer: %w", err)
	}

	// Check if GigaChat wants to call a function
	if result.FinishReason == "function_call" && result.FunctionCall != nil {
		// Call MCP tool
		toolResult, err := a.callMCPTool(result.FunctionCall)
		if err != nil {
			return "", fmt.Errorf("failed to call MCP tool: %w", err)
		}

		// Add assistant's function call to messages
		a.messages = append(a.messages, *result.Message)

		// Add function result as user message
		a.messages = append(a.messages, gigachat.Message{
			Role:    "user",
			Content: fmt.Sprintf("Function %s result: %s", result.FunctionCall.Name, toolResult),
		})

		// Get final answer from GigaChat
		finalResult, err := a.gigaChatNetworkService.GetCompletion(a.messages, a.Model, a.temperature, a.functions)
		if err != nil {
			return "", fmt.Errorf("failed to get final answer: %w", err)
		}

		a.messages = append(a.messages, *finalResult.Message)
		return finalResult.Message.Content, nil
	}

	// No function call, just return the message
	a.messages = append(a.messages, *result.Message)
	return result.Message.Content, nil
}

func (a *Agent) callMCPTool(functionCall *gigachat.FunctionCall) (string, error) {
	// Arguments are already parsed as map[string]any
	args := functionCall.Arguments

	// Call MCP tool
	ctx := context.Background()
	result, err := a.mcpClient.CallTool(ctx, functionCall.Name, args)
	if err != nil {
		return "", err
	}

	// Extract result content
	if len(result.Content) == 0 {
		return "", fmt.Errorf("no content in tool result")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		return "", fmt.Errorf("unexpected content type")
	}

	return textContent.Text, nil
}