package xai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// InvokeStream sends a streaming chat completion request to the xAI API
func (c *Client) InvokeStream(
	messages []map[string]interface{},
	opts ...InvokeOption,
) (<-chan *ChatCompletionResponse, error) {
	// Create request with streaming enabled
	streamOpt := WithStream(true)
	allOpts := append([]InvokeOption{streamOpt}, opts...)

	req := &ChatCompletionRequest{
		Messages: messages,
		Model:    c.Model,
	}

	// Apply options
	for _, opt := range allOpts {
		opt(req)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if tools are provided
	if len(req.Tools) > 0 {
		return nil, fmt.Errorf("tools are not supported in streaming mode")
	}

	// Make streaming API call
	httpResp, err := c.makeStreamingAPICall(req)
	if err != nil {
		return nil, err
	}

	// Create channel for streaming responses
	respChan := make(chan *ChatCompletionResponse)

	// Process stream in a goroutine
	go func() {
		defer httpResp.Body.Close()
		defer close(respChan)

		scanner := bufio.NewScanner(httpResp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			line = strings.TrimPrefix(line, "data: ")
			if line == "[DONE]" {
				break
			}

			var chunk ChatCompletionResponse
			if err := json.Unmarshal([]byte(line), &chunk); err != nil {
				// Skip malformed chunks
				continue
			}

			respChan <- &chunk
		}

		if err := scanner.Err(); err != nil {
			// Send error through channel or log it
			// For now, we just let the channel close
		}
	}()

	return respChan, nil
}

// makeStreamingAPICall sends a streaming request to the xAI API
func (c *Client) makeStreamingAPICall(req *ChatCompletionRequest) (*http.Response, error) {
	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body[:n]))
	}

	return resp, nil
}
