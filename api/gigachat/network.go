package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	oauthURL       = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	completionsURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	rqUID          = "270fee8f-3594-4cb7-b9cb-d0690691f735"
)

type NetworkService struct{
	oauthToken string
	requestToken string
}

func GetNetworkService() (*NetworkService, error) {
	networkService := &NetworkService{
		oauthToken:		GetOauthToken(),
		requestToken:	"",
	}

	requestToken, err := networkService.GetRequestToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get request token: %w", err)
	}
	networkService.requestToken = requestToken

	return networkService, nil
}

func (networkService *NetworkService) GetRequestToken() (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequest("POST", oauthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //what does it mean?
	req.Header.Set("Accept", "application/json") //what does it mean? Accept
	req.Header.Set("Authorization", "Bearer "+networkService.oauthToken)
	req.Header.Set("RqUID", rqUID)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (networkService *NetworkService) GetCompletion(messages []Message, model string, temperature float64, functions []Function) (*CompletionResult, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	reqData := CompletionRequest{
		Model:             model,
		Messages:          messages,
		Temperature:       temperature,
		RepetitionPenalty: 1,
	}

	// Add functions if provided
	if len(functions) > 0 {
		reqData.FunctionCall = "auto"
		reqData.Functions = functions
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", completionsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+networkService.requestToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var completionResp CompletionResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(completionResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	choice := completionResp.Choices[0]

	return &CompletionResult{
		Message: &Message{
			Role:         choice.Message.Role,
			Content:      choice.Message.Content,
			FunctionCall: choice.Message.FunctionCall,
		},
		FinishReason: choice.FinishReason,
		FunctionCall: choice.Message.FunctionCall,
	}, nil
}
