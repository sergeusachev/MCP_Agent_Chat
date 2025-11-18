package network

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client with configurable settings
type Client struct {
	httpClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new network client with default settings
func NewClient(timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// Request represents an HTTP request configuration
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Do executes an HTTP request and returns the response
func (c *Client) Do(req *Request) (*Response, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute the request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &Response{
		StatusCode: httpResp.StatusCode,
		Body:       body,
		Headers:    httpResp.Header,
	}

	// Check for HTTP errors
	if httpResp.StatusCode >= 400 {
		return response, fmt.Errorf("HTTP error: status code %d, body: %s", httpResp.StatusCode, string(body))
	}

	return response, nil
}

// Get performs a GET request
func (c *Client) Get(url string, headers map[string]string) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
	})
}

// Post performs a POST request
func (c *Client) Post(url string, headers map[string]string, body io.Reader) (*Response, error) {
	return c.Do(&Request{
		Method:  http.MethodPost,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}
