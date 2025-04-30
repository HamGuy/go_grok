package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hamguy/go_grok/pkg/utils"
	"github.com/hamguy/go_grok/pkg/xai"
)

func main() {
	// Get API key
	apiKey := utils.GetAPIKey()
	if apiKey == "" {
		log.Fatal("No API key found, please set the GROK_API_KEY environment variable or .env file")
	}

	fmt.Printf("Using API key: %s\n", utils.GetMaskedAPIKey(apiKey))

	// Initialize client
	client := xai.NewClient(
		apiKey,
		string(xai.Grok3Beta), // Using Grok3 model
	)

	// Simulate loading animation
	done := make(chan bool)
	go func() {
		fmt.Print("Getting response")
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				fmt.Println()
				return
			case <-ticker.C:
				fmt.Print(".")
			}
		}
	}()

	// Convert messages to map format
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Please tell me who you are, and answer in Chinese",
		},
	}

	// Set temperature parameter
	temp := 0.7

	// Create request
	req := &xai.ChatCompletionRequest{
		Messages:    messages,
		Temperature: &temp,
	}

	// Streaming request
	response, err := client.CreateChatCompletionStream(req)
	close(done) // Stop loading animation

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Process streaming response
	var fullContent strings.Builder
	for chunk := range response.Stream {
		content := chunk.Choices[0].Delta.Content
		fmt.Print(content)
		fullContent.WriteString(content)
	}

	fmt.Println("\n\nFull answer:")
	fmt.Println(fullContent.String())

	// English example
	fmt.Println("\nEnglish Example:")
	done = make(chan bool)
	go func() {
		fmt.Print("Getting response")
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				fmt.Println()
				return
			case <-ticker.C:
				fmt.Print(".")
			}
		}
	}()

	// English message
	englishMessages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Tell me who you are, and answer in English",
		},
	}

	// Create English request
	englishReq := &xai.ChatCompletionRequest{
		Messages:    englishMessages,
		Temperature: &temp,
	}

	response, err = client.CreateChatCompletionStream(englishReq)
	close(done)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fullContent.Reset()
	for chunk := range response.Stream {
		content := chunk.Choices[0].Delta.Content
		fmt.Print(content)
		fullContent.WriteString(content)
	}

	fmt.Println("\n\nFull answer:")
	fmt.Println(fullContent.String())
}
