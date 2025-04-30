package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hamguy/go_grok/pkg/xai"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set. Please set it and try again.")
	}

	// Create a new client
	client := xai.NewClient(
		apiKey,
		string(xai.Grok3Beta), // Using Grok 3 model
	)

	// Create a simple message
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Hello! What can you tell me about yourself?",
		},
	}

	// Set options for the request
	options := []xai.InvokeOption{
		xai.WithTemperature(0.7),
		xai.WithMaxTokens(500),
	}

	// Make the API call
	resp, err := client.Invoke(messages, options...)
	if err != nil {
		log.Fatalf("Error calling API: %v", err)
	}

	// Extract and print the response
	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
		fmt.Println("Response:")
		fmt.Println(resp.Choices[0].Message.Content)
	} else {
		fmt.Println("No response content received")
	}

	// Print usage statistics
	if resp.Usage != nil {
		fmt.Printf("\nToken usage:\n")
		fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}
