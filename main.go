package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hamguy/go_grok/pkg/xai"
)

// printMaskedKey prints a masked version of the API key
func printMaskedKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	// Only show the first 4 and last 4 characters, replace the middle with ****
	return key[:4] + "****" + key[len(key)-4:]
}

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		// Try to read the key from file
		log.Printf("GROK_API_KEY not found in environment variables, trying to read from .env file")
		keyBytes, err := os.ReadFile("../../.env")
		if err == nil {
			log.Printf("Found .env file, parsing content")
			lines := strings.Split(string(keyBytes), "\n")
			for _, line := range lines {
				// Debug output for each line (masked)
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
					log.Printf("Processing line: %s", trimmedLine[:3]+"...")
				}

				if strings.HasPrefix(trimmedLine, "GROK_API_KEY=") {
					apiKey = strings.TrimPrefix(trimmedLine, "GROK_API_KEY=")
					log.Printf("Found API key in .env file")
					break
				}
			}
		} else {
			log.Printf("Error reading .env file: %v", err)
		}
	}

	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set. Please set it and try again.")
	}

	// Print masked version of API key
	log.Printf("Using API key: %s", printMaskedKey(apiKey))

	// Create a new client
	client := xai.NewClient(
		apiKey,
		string(xai.Grok3Beta), // Using Grok 3 model
	)

	// Create a message that will generate a longer response
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Write a short poem about artificial intelligence.",
		},
	}

	fmt.Println("Streaming response from Grok AI:")
	fmt.Println("--------------------------------")

	// Start streaming the response
	streamChan, err := client.InvokeStream(messages)
	if err != nil {
		log.Fatalf("Error calling InvokeStream: %v", err)
	}

	// Simple animation to show we're waiting for the first chunk
	doneChan := make(chan bool)
	go func() {
		chars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
		i := 0
		for {
			select {
			case <-doneChan:
				return
			default:
				fmt.Printf("\rWaiting for response %s", chars[i])
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Track if we've received our first chunk
	receivedChunk := false

	// Process and display the streaming response
	for chunk := range streamChan {
		if !receivedChunk {
			// Clear the loading animation line
			fmt.Print("\r                          \r")
			receivedChunk = true
			close(doneChan)
		}

		// Print the content of each chunk
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	fmt.Println("\n--------------------------------")
	fmt.Println("Stream complete!")
}
