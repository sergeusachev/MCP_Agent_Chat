package gigachat

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"` // JSON schema as object or string
}

type FunctionCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"` // Function arguments as object
}

type CompletionRequest struct {
	Model             string     `json:"model"`
	Messages          []Message  `json:"messages"`
	FunctionCall      string     `json:"function_call,omitempty"`
	Functions         []Function `json:"functions,omitempty"`
	Temperature       float64    `json:"temperature"`
	RepetitionPenalty float64    `json:"repetition_penalty"`
}

type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Role         string        `json:"role"`
			Content      string        `json:"content"`
			FunctionCall *FunctionCall `json:"function_call,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type CompletionResult struct {
	Message       *Message
	FinishReason  string
	FunctionCall  *FunctionCall
}