package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hamguy/go_grok/pkg/xai"
)

func main() {
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	client := xai.NewClient(
		apiKey,
		string(xai.Grok212),
	)

	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "Explain quantum computing in simple terms.",
		},
	}

	// Start a goroutine to show a spinner while waiting for the first chunk
	spinnerDone := make(chan bool)
	go func() {
		spinners := []string{"|", "/", "-", "\\"}
		i := 0
		for {
			select {
			case <-spinnerDone:
				return
			default:
				fmt.Printf("\rThinking %s", spinners[i])
				i = (i + 1) % len(spinners)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Call the streaming API
	streamChan, err := client.InvokeStream(messages)
	if err != nil {
		close(spinnerDone)
		log.Fatalf("Error calling InvokeStream: %v", err)
	}

	// Variables to track streaming progress
	var fullResponse strings.Builder
	receivedFirstChunk := false

	// Process the streaming response
	for chunk := range streamChan {
		if !receivedFirstChunk {
			// Clear the spinner line when we get the first chunk
			fmt.Print("\r                      \r")
			receivedFirstChunk = true
			close(spinnerDone)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content := chunk.Choices[0].Delta.Content
			fullResponse.WriteString(content)
			fmt.Print(content)
			os.Stdout.Sync()
		}
	}

	fmt.Println("\n\nFull response received!")
}
