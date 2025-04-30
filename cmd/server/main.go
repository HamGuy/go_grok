package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hamguy/go_grok/pkg/xai"
)

// Server configuration
const (
	defaultPort = "8080"
	defaultHost = "localhost"
)

// Request and response structures
type ChatRequest struct {
	Messages []map[string]interface{} `json:"messages"`
	Model    string                   `json:"model,omitempty"`
	Stream   bool                     `json:"stream,omitempty"`
	// Other optional parameters
	Temperature *float64                 `json:"temperature,omitempty"`
	MaxTokens   *int                     `json:"max_tokens,omitempty"`
	Tools       []map[string]interface{} `json:"tools,omitempty"`
	ToolChoice  interface{}              `json:"tool_choice,omitempty"`
}

func main() {
	// Get API Key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set. Please set it and try again.")
	}

	// Determine server port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create global client
	client := xai.NewClient(
		apiKey,
		string(xai.Grok3Beta),
	)

	// Define handler function
	http.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse request
		var chatReq ChatRequest
		if err := json.Unmarshal(body, &chatReq); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Prepare options
		var options []xai.InvokeOption
		if chatReq.Temperature != nil {
			options = append(options, xai.WithTemperature(*chatReq.Temperature))
		}
		if chatReq.MaxTokens != nil {
			options = append(options, xai.WithMaxTokens(*chatReq.MaxTokens))
		}

		// Set up the client with tools if provided
		if len(chatReq.Tools) > 0 {
			// Since WithTools is a ClientOption, we need to recreate the client
			client = xai.NewClient(
				apiKey,
				string(xai.Grok3Beta),
				xai.WithTools(chatReq.Tools),
			)
		}

		// Add tool choice option if provided
		if chatReq.ToolChoice != nil {
			options = append(options, xai.WithToolChoice(chatReq.ToolChoice))
		}

		// Handle streaming response
		if chatReq.Stream {
			handleStreamingResponse(w, client, chatReq.Messages, options)
			return
		}

		// Handle standard response
		handleStandardResponse(w, client, chatReq.Messages, options)
	})

	// Add health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", defaultHost, port)
	log.Printf("🚀 Server started at http://%s", serverAddr)
	log.Printf("💡 Example request: curl -X POST http://%s/chat/completions -H 'Content-Type: application/json' -d '{\"messages\":[{\"role\":\"user\",\"content\":\"Hello!\"}]}'", serverAddr)

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}

// Handle standard response
func handleStandardResponse(w http.ResponseWriter, client *xai.Client, messages []map[string]interface{}, options []xai.InvokeOption) {
	// Call API
	resp, err := client.Invoke(messages, options...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling Grok API: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Handle streaming response
func handleStreamingResponse(w http.ResponseWriter, client *xai.Client, messages []map[string]interface{}, options []xai.InvokeOption) {
	// Set response headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	// Flush writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Get streaming response
	streamChan, err := client.InvokeStream(messages, options...)
	if err != nil {
		// Send error in streaming response
		fmt.Fprintf(w, "data: {\"error\": \"%v\"}\n\n", err)
		flusher.Flush()
		return
	}

	// Process streaming response
	for chunk := range streamChan {
		// Convert response chunk to JSON
		chunkBytes, err := json.Marshal(chunk)
		if err != nil {
			continue
		}

		// Send response chunk
		fmt.Fprintf(w, "data: %s\n\n", chunkBytes)
		flusher.Flush()
	}

	// Send end marker
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
