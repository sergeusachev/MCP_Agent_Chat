package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"serge.com/mcp-example/common/network"
)

const (
	oauthURL       = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	completionsURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	defaultTimeout = 30 * time.Second
)

// GigaChatClient represents a client for the GigaChat API
type GigaChatClient struct {
	networkClient *network.Client
	oauthToken    string
	accessToken   string
}

// OAuthResponse represents the response from the OAuth endpoint
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest represents a request to the completions endpoint
type CompletionRequest struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Temperature       float64   `json:"temperature,omitempty"`
	RepetitionPenalty float64   `json:"repetition_penalty,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Message      Message `json:"message"`
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CompletionResponse represents the response from the completions endpoint
type CompletionResponse struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Object  string   `json:"object"`
}

// NewGigaChatClient creates a new GigaChat API client
func NewGigaChatClient() (*GigaChatClient, error) {
	oauthToken, err := loadOAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to load OAuth token: %w", err)
	}

	// Create network client with custom transport that skips SSL verification
	// (GigaChat API uses self-signed certificates)
	httpClient := &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client := &GigaChatClient{
		networkClient: &network.Client{},
		oauthToken:    oauthToken,
	}

	// We need to use custom HTTP client for GigaChat, so let's create it here
	// Get access token on initialization
	accessToken, err := client.getAccessToken(httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	client.accessToken = accessToken

	return client, nil
}

// loadOAuthToken loads the OAuth token from the secret/oauth_gigachat_token.txt file
func loadOAuthToken() (string, error) {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Build path to secret/oauth_gigachat_token.txt relative to this source file
	packageDir := filepath.Dir(filename)
	tokenPath := filepath.Join(packageDir, "secret", "oauth_gigachat_token.txt")

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read OAuth token file: %w", err)
	}

	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("OAuth token is empty")
	}

	return token, nil
}

// getAccessToken gets an access token from the OAuth endpoint
func (c *GigaChatClient) getAccessToken(httpClient *http.Client) (string, error) {
	// Generate a random RqUID (can be any UUID)
	rqUID := "270fee8f-3594-4cb7-b9cb-d0690691f735"

	// Prepare request body
	body := strings.NewReader("scope=GIGACHAT_API_PERS")

	// Create HTTP request
	req, err := http.NewRequest("POST", oauthURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.oauthToken)
	req.Header.Set("RqUID", rqUID)

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OAuth request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var oauthResp OAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return "", fmt.Errorf("failed to parse OAuth response: %w", err)
	}

	if oauthResp.AccessToken == "" {
		return "", fmt.Errorf("access token is empty in response")
	}

	return oauthResp.AccessToken, nil
}

// SendCompletion sends a chat completion request to GigaChat
func (c *GigaChatClient) SendCompletion(req *CompletionRequest) (*CompletionResponse, error) {
	// Create HTTP client with SSL skip
	httpClient := &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Marshal request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", completionsURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.accessToken)

	// Execute request
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("completions request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var completionResp CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return nil, fmt.Errorf("failed to parse completion response: %w", err)
	}

	return &completionResp, nil
}

// Chat is a convenience method to send a simple chat message
func (c *GigaChatClient) Chat(userMessage string) (string, error) {
	req := &CompletionRequest{
		Model: "GigaChat-2",
		Messages: []Message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		Temperature:       0,
		RepetitionPenalty: 1,
	}

	resp, err := c.SendCompletion(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}
