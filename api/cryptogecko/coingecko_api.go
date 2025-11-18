package cryptogecko

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"serge.com/mcp-example/common/network"
)

const (
	baseURL        = "https://api.coingecko.com/api/v3"
	apiKeyHeader   = "x-cg-demo-api-key"
	defaultTimeout = 30 * time.Second
)

// CoinGeckoClient represents a client for the CoinGecko API
type CoinGeckoClient struct {
	networkClient *network.Client
	apiKey        string
}

// PriceResponse represents a cryptocurrency price
type PriceResponse struct {
	CoinID   string
	Currency string
	Price    float64
}

// NewCoinGeckoClient creates a new CoinGecko API client
func NewCoinGeckoClient() (*CoinGeckoClient, error) {
	apiKey, err := loadAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load API key: %w", err)
	}

	return &CoinGeckoClient{
		networkClient: network.NewClient(defaultTimeout),
		apiKey:        apiKey,
	}, nil
}

// loadAPIKey loads the API key from the secret/api_key.txt file
func loadAPIKey() (string, error) {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Build path to secret/api_key.txt relative to this source file
	packageDir := filepath.Dir(filename)
	apiKeyPath := filepath.Join(packageDir, "secret", "api_key.txt")

	data, err := os.ReadFile(apiKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key file: %w", err)
	}

	apiKey := strings.TrimSpace(string(data))
	if apiKey == "" {
		return "", fmt.Errorf("API key is empty")
	}

	return apiKey, nil
}

// GetCoinPrice gets the price of a specific coin in a specific currency
func (c *CoinGeckoClient) GetCoinPrice(coinID string, currency string) (*PriceResponse, error) {
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", baseURL, coinID, currency)

	headers := map[string]string{
		apiKeyHeader: c.apiKey,
	}

	resp, err := c.networkClient.Get(url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price: %w", err)
	}

	// Parse the nested map response from API
	var apiResponse map[string]map[string]float64
	if err := json.Unmarshal(resp.Body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract the price from the nested map
	coinData, ok := apiResponse[coinID]
	if !ok {
		return nil, fmt.Errorf("coin %s not found in response", coinID)
	}

	price, ok := coinData[currency]
	if !ok {
		return nil, fmt.Errorf("currency %s not found for coin %s", currency, coinID)
	}

	return &PriceResponse{
		CoinID:   coinID,
		Currency: currency,
		Price:    price,
	}, nil
}
