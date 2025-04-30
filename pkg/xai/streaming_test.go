// pkg/xai/streaming_test.go

package xai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestStreamingOutput tests the streaming output functionality of the SDK
func TestStreamingOutput(t *testing.T) {
	// Define a constant for the API key to avoid hardcoding
	const testAPIKey = "your_api_key_here"

	// Setup a mock server that simulates streaming responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request is properly formed
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected path /chat/completions, got %s", r.URL.Path)
		}

		// Verify headers
		if r.Header.Get("Authorization") != "Bearer "+testAPIKey {
			t.Errorf("Expected Authorization header 'Bearer %s', got %s", testAPIKey, r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("Expected Accept header 'text/event-stream', got %s", r.Header.Get("Accept"))
		}

		// Set headers for streaming response
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Create a flusher to ensure each chunk is sent immediately
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("ResponseWriter does not implement http.Flusher")
			return
		}

		// Define the streaming chunks to send
		// Each chunk represents a part of the response
		chunks := []string{
			`data: {"id":"stream-test-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"role":"assistant","content":"This"},"finish_reason":null}]}`,
			`data: {"id":"stream-test-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":" is"},"finish_reason":null}]}`,
			`data: {"id":"stream-test-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":" a"},"finish_reason":null}]}`,
			`data: {"id":"stream-test-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":" streaming"},"finish_reason":null}]}`,
			`data: {"id":"stream-test-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":" test"},"finish_reason":null}]}`,
			`data: {"id":"stream-test-id","object":"chat.completion.chunk","created":1677858242,"model":"grok-2-1212","choices":[{"index":0,"delta":{"content":"."},"finish_reason":"stop"}]}`,
			`data: [DONE]`,
		}

		// Send each chunk with a small delay to simulate real streaming
		for _, chunk := range chunks {
			fmt.Fprintf(w, "%s\n\n", chunk)
			flusher.Flush()
			time.Sleep(50 * time.Millisecond) // Small delay between chunks
		}
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewClient(
		testAPIKey,
		string(Grok212),
		WithBaseURL(server.URL),
	)

	// Create a test message
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Test streaming response",
		},
	}

	// Call the streaming API
	streamChan, err := client.InvokeStream(messages)
	if err != nil {
		t.Fatalf("InvokeStream returned an error: %v", err)
	}

	// Collect the streaming response
	var fullContent string
	var chunks []string

	for chunk := range streamChan {
		// Verify each chunk has the expected structure
		if chunk.ID != "stream-test-id" {
			t.Errorf("Expected chunk ID 'stream-test-id', got '%s'", chunk.ID)
		}

		if len(chunk.Choices) != 1 {
			t.Errorf("Expected 1 choice in chunk, got %d", len(chunk.Choices))
		}

		// Extract content from the delta
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content := chunk.Choices[0].Delta.Content
			chunks = append(chunks, content)
			fullContent += content
		}
	}

	// Verify the complete content
	expectedContent := "This is a streaming test."
	if fullContent != expectedContent {
		t.Errorf("Expected full content '%s', got '%s'", expectedContent, fullContent)
	}

	// Verify we received the expected number of chunks
	expectedChunks := []string{"This", " is", " a", " streaming", " test", "."}
	if len(chunks) != len(expectedChunks) {
		t.Errorf("Expected %d chunks, got %d", len(expectedChunks), len(chunks))
	}

	// Verify each chunk's content
	for i, chunk := range chunks {
		if i < len(expectedChunks) && chunk != expectedChunks[i] {
			t.Errorf("Chunk %d: expected '%s', got '%s'", i, expectedChunks[i], chunk)
		}
	}
}
