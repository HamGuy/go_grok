package xai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is the main client for interacting with the xAI API
type Client struct {
	APIKey      string
	Model       string
	BaseURL     string
	Tools       []map[string]interface{}
	FunctionMap map[string]func(map[string]interface{}) (string, error)
}

// NewClient creates a new xAI client
func NewClient(apiKey string, model string, options ...ClientOption) *Client {
	client := &Client{
		APIKey:      apiKey,
		Model:       model,
		BaseURL:     "https://api.x.ai/v1",
		Tools:       []map[string]interface{}{},
		FunctionMap: make(map[string]func(map[string]interface{}) (string, error)),
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = url
	}
}

// WithTools sets the available tools for the client
func WithTools(tools []map[string]interface{}) ClientOption {
	return func(c *Client) {
		c.Tools = tools
	}
}

// WithFunctionMap sets the function map for tool execution
func WithFunctionMap(functionMap map[string]func(map[string]interface{}) (string, error)) ClientOption {
	return func(c *Client) {
		c.FunctionMap = functionMap
	}
}

// Invoke sends a chat completion request to the xAI API
func (c *Client) Invoke(
	messages []map[string]interface{},
	opts ...InvokeOption,
) (*ChatCompletionResponse, error) {
	// Create request with default values
	req := &ChatCompletionRequest{
		Messages: messages,
		Model:    c.Model,
		Tools:    c.Tools,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check for streaming
	if req.Stream != nil && *req.Stream {
		return nil, fmt.Errorf("for streaming, use InvokeStream instead of Invoke")
	}

	// Make API call
	resp, err := c.makeAPICall(req)
	if err != nil {
		return nil, err
	}

	// Process tool calls if present
	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil && len(resp.Choices[0].Message.ToolCalls) > 0 {
		c.processToolCalls(resp)
	}

	return resp, nil
}

// makeAPICall sends a request to the xAI API
func (c *Client) makeAPICall(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var chatResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}

// processToolCalls executes tool calls and adds results to the response
func (c *Client) processToolCalls(resp *ChatCompletionResponse) {
	if len(resp.Choices) == 0 || resp.Choices[0].Message == nil {
		return
	}

	message := resp.Choices[0].Message
	toolResults := []ToolResult{}

	for _, toolCall := range message.ToolCalls {
		if toolCall.Type == "function" {
			funcName := toolCall.Function.Name
			if fn, ok := c.FunctionMap[funcName]; ok {
				// Parse arguments - handle both string and map types
				var args map[string]interface{}

				switch v := toolCall.Function.Arguments.(type) {
				case string:
					// If it's a JSON string, parse it
					if err := json.Unmarshal([]byte(v), &args); err != nil {
						toolResults = append(toolResults, ToolResult{
							ToolCallID: toolCall.ID,
							Role:       "tool",
							Name:       funcName,
							Content:    fmt.Sprintf("Error parsing arguments: %s", err.Error()),
						})
						continue
					}
				case map[string]interface{}:
					// If it's already a map, use it directly
					args = v
				default:
					// If it's neither a string nor a map, report an error
					toolResults = append(toolResults, ToolResult{
						ToolCallID: toolCall.ID,
						Role:       "tool",
						Name:       funcName,
						Content:    fmt.Sprintf("Error: unexpected arguments type %T", v),
					})
					continue
				}

				// Execute function
				result, err := fn(args)
				if err != nil {
					result = fmt.Sprintf("Error: %s", err.Error())
				}

				// Add result
				toolResults = append(toolResults, ToolResult{
					ToolCallID: toolCall.ID,
					Role:       "tool",
					Name:       funcName,
					Content:    result,
				})
			}
		}
	}

	// Add tool results to message
	if len(toolResults) > 0 {
		message.ToolResults = toolResults
	}
}
