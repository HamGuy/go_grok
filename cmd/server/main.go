package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hamguy/xai_grok_sdk_go/pkg/xai"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Define a tool function for getting weather
	getWeather := func(args map[string]interface{}) (string, error) {
		location, ok := args["location"].(string)
		if !ok {
			return "", fmt.Errorf("location must be a string")
		}
		// In a real implementation, you would call a weather API here
		return fmt.Sprintf("The weather in %s is sunny and 72°F.", location), nil
	}

	// Define tools
	tools := []map[string]interface{}{
		{
			"name":        "get_weather",
			"description": "Get the current weather for a location",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city and state, e.g., San Francisco, CA",
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

	// Initialize client
	client := xai.NewClient(
		apiKey,
		string(xai.Grok212),
		xai.WithTools(tools),
		xai.WithFunctionMap(functionMap),
	)

	// Example 1: Basic completion
	fmt.Println("Example 1: Basic completion")
	resp, err := client.Invoke(
		[]map[string]interface{}{
			{"role": "user", "content": "Hello, how are you today?"},
		},
		xai.WithTemperature(0.7),
		xai.WithMaxTokens(100),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
		fmt.Printf("Response: %s\n\n", resp.Choices[0].Message.Content)
	}

	// Example 2: Function calling
	fmt.Println("Example 2: Function calling")
	resp, err = client.Invoke(
		[]map[string]interface{}{
			{"role": "user", "content": "What's the weather like in San Francisco?"},
		},
		xai.WithToolChoice("auto"),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
		message := resp.Choices[0].Message
		fmt.Printf("Assistant: %s\n", message.Content)

		if len(message.ToolResults) > 0 {
			fmt.Printf("Tool Result: %s\n\n", message.ToolResults[0].Content)
		}
	}

	// Example 3: Streaming
	fmt.Println("Example 3: Streaming")
	streamChan, err := client.InvokeStream(
		[]map[string]interface{}{
			{"role": "user", "content": "Write a short poem about artificial intelligence."},
		},
		xai.WithTemperature(0.8),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Print("Streaming response: ")
	for chunk := range streamChan {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println("\n")
}
