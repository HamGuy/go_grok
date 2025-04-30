// pkg/xai/client_test.go

package xai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAPIKey = "your_api_key_here"

// TestClientInvoke tests the basic functionality of the Client.Invoke method
func TestClientInvoke(t *testing.T) {
	// Setup a mock server to simulate the xAI API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected path /chat/completions, got %s", r.URL.Path)
		}

		// Verify authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testAPIKey {
			t.Errorf("Expected Authorization header 'Bearer %s', got %s", testAPIKey, authHeader)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got %s", contentType)
		}

		// Decode the request body to verify parameters
		var req ChatCompletionRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify request parameters
		if req.Model != "grok-2-1212" {
			t.Errorf("Expected model 'grok-2-1212', got %s", req.Model)
		}
		if len(req.Messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(req.Messages))
		}
		if req.Messages[0]["role"] != "user" {
			t.Errorf("Expected role 'user', got %s", req.Messages[0]["role"])
		}
		if req.Messages[0]["content"] != "Hello, world!" {
			t.Errorf("Expected content 'Hello, world!', got %s", req.Messages[0]["content"])
		}

		// Return a mock response
		mockResponse := `{  
			"id": "test-id",  
			"object": "chat.completion",  
			"created": 1677858242,  
			"model": "grok-2-1212",  
			"choices": [  
				{  
					"index": 0,  
					"message": {  
						"role": "assistant",  
						"content": "Hello! How can I help you today?"  
					},  
					"finish_reason": "stop"  
				}  
			],  
			"usage": {  
				"prompt_tokens": 10,  
				"completion_tokens": 8,  
				"total_tokens": 18  
			},  
			"system_fingerprint": "test-fingerprint"  
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create a client that uses our mock server
	client := NewClient(
		testAPIKey,
		string(Grok212),
		WithBaseURL(server.URL), // Use the mock server URL
	)

	// Test the Invoke method
	resp, err := client.Invoke(
		[]map[string]interface{}{
			{"role": "user", "content": "Hello, world!"},
		},
		WithTemperature(0.7),
	)

	// Check for errors
	if err != nil {
		t.Fatalf("Invoke returned an error: %v", err)
	}

	// Verify response
	if resp.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", resp.ID)
	}
	if resp.Model != "grok-2-1212" {
		t.Errorf("Expected model 'grok-2-1212', got %s", resp.Model)
	}
	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got %s", resp.Choices[0].Message.Role)
	}
	if resp.Choices[0].Message.Content != "Hello! How can I help you today?" {
		t.Errorf("Expected content 'Hello! How can I help you today?', got %s", resp.Choices[0].Message.Content)
	}
	if resp.Usage == nil {
		t.Error("Expected usage to be non-nil")
	} else {
		if resp.Usage.PromptTokens != 10 {
			t.Errorf("Expected prompt_tokens 10, got %d", resp.Usage.PromptTokens)
		}
		if resp.Usage.CompletionTokens != 8 {
			t.Errorf("Expected completion_tokens 8, got %d", resp.Usage.CompletionTokens)
		}
		if resp.Usage.TotalTokens != 18 {
			t.Errorf("Expected total_tokens 18, got %d", resp.Usage.TotalTokens)
		}
	}
}

// TestClientToolCalling tests the tool calling functionality
func TestClientToolCalling(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a mock response with tool calls
		mockResponse := `{  
			"id": "test-tool-id",  
			"object": "chat.completion",  
			"created": 1677858242,  
			"model": "grok-2-1212",  
			"choices": [  
				{  
					"index": 0,  
					"message": {  
						"role": "assistant",  
						"content": "I'll get the weather for you.",  
						"tool_calls": [  
							{  
								"id": "call_123",  
								"type": "function",  
								"function": {  
									"name": "get_weather",  
									"arguments": "{\"location\":\"San Francisco\"}"  
								}  
							}  
						]  
					},  
					"finish_reason": "tool_calls"  
				}  
			],  
			"usage": {  
				"prompt_tokens": 15,  
				"completion_tokens": 12,  
				"total_tokens": 27  
			},  
			"system_fingerprint": "test-fingerprint"  
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Define a mock weather function
	getWeather := func(args map[string]interface{}) (string, error) {
		location, _ := args["location"].(string)
		return "The weather in " + location + " is sunny.", nil
	}

	// Define tools
	tools := []map[string]interface{}{
		{
			"name":        "get_weather",
			"description": "Get the weather for a location",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city name",
					},
				},
				"required": []string{"location"},
			},
		},
	}

	// Create function map
	functionMap := map[string]func(map[string]interface{}) (string, error){
		"get_weather": getWeather,
	}

	// Create a client
	client := NewClient(
		testAPIKey,
		string(Grok212),
		WithBaseURL(server.URL),
		WithTools(tools),
		WithFunctionMap(functionMap),
	)

	// Test the Invoke method with tool calling
	resp, err := client.Invoke(
		[]map[string]interface{}{
			{"role": "user", "content": "What's the weather in San Francisco?"},
		},
		WithToolChoice("auto"),
	)

	// Check for errors
	if err != nil {
		t.Fatalf("Invoke returned an error: %v", err)
	}

	// Verify tool calls in response
	if len(resp.Choices[0].Message.ToolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(resp.Choices[0].Message.ToolCalls))
	} else {
		toolCall := resp.Choices[0].Message.ToolCalls[0]
		if toolCall.Type != "function" {
			t.Errorf("Expected tool call type 'function', got %s", toolCall.Type)
		}
		if toolCall.Function.Name != "get_weather" {
			t.Errorf("Expected function name 'get_weather', got %s", toolCall.Function.Name)
		}
	}

	// Verify tool results
	if len(resp.Choices[0].Message.ToolResults) != 1 {
		t.Errorf("Expected 1 tool result, got %d", len(resp.Choices[0].Message.ToolResults))
	} else {
		toolResult := resp.Choices[0].Message.ToolResults[0]
		if toolResult.Name != "get_weather" {
			t.Errorf("Expected tool result name 'get_weather', got %s", toolResult.Name)
		}
		if toolResult.Content != "The weather in San Francisco is sunny." {
			t.Errorf("Expected content 'The weather in San Francisco is sunny.', got %s", toolResult.Content)
		}
	}
}

// TestClientStreaming tests the streaming functionality
func TestClientStreaming(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify stream parameter
		var req ChatCompletionRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Stream == nil || !*req.Stream {
			t.Errorf("Expected stream to be true")
		}

		// Set headers for streaming
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Send streaming chunks
		chunks := []string{
			`data: {"id":"test-stream-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`,
			`data: {"id":"test-stream-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":" there"},"finish_reason":null}]}`,
			`data: {"id":"test-stream-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}`,
			`data: [DONE]`,
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("Expected ResponseWriter to be a Flusher")
			return
		}

		for _, chunk := range chunks {
			fmt.Fprintf(w, "%s\n\n", chunk)
			flusher.Flush()
		}
	}))
	defer server.Close()

	// Create a client
	client := NewClient(
		testAPIKey,
		string(Grok212),
		WithBaseURL(server.URL),
	)

	// Test the InvokeStream method
	streamChan, err := client.InvokeStream(
		[]map[string]interface{}{
			{"role": "user", "content": "Say hello"},
		},
	)

	// Check for errors
	if err != nil {
		t.Fatalf("InvokeStream returned an error: %v", err)
	}

	// Collect all chunks
	var content string
	for chunk := range streamChan {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content += chunk.Choices[0].Delta.Content
		}
	}

	// Verify the combined content
	expectedContent := "Hello there!"
	if content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, content)
	}
}
