package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hamguy/go_grok/pkg/xai"
)

// getCurrentWeather returns weather information for a location
func getCurrentWeather(args map[string]interface{}) (string, error) {
	location, ok := args["location"].(string)
	if !ok {
		return "", fmt.Errorf("location must be a string")
	}

	// In a real application, you would call a weather API here
	// For the example, we'll return mock data
	return fmt.Sprintf("Current weather in %s: 22°C, Sunny with light clouds", location), nil
}

// getStockPrice returns stock price information
func getStockPrice(args map[string]interface{}) (string, error) {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return "", fmt.Errorf("symbol must be a string")
	}

	// In a real application, you would call a financial API here
	// For the example, we'll return mock data with the current timestamp
	price := 150.25 + float64(time.Now().UnixNano()%1000)/100
	return fmt.Sprintf("%s is currently trading at $%.2f", symbol, price), nil
}

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set. Please set it and try again.")
	}

	// Define tools that the model can use
	tools := []map[string]interface{}{
		{
			"name":        "get_current_weather",
			"description": "Get the current weather in a given location",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city name, e.g., Beijing, Shanghai",
					},
				},
				"required": []string{"location"},
			},
		},
		{
			"name":        "get_stock_price",
			"description": "Get the current stock price for a given symbol",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "The stock symbol, e.g., AAPL for Apple",
					},
				},
				"required": []string{"symbol"},
			},
		},
	}

	// Map tool names to their handler functions
	functionMap := map[string]func(map[string]interface{}) (string, error){
		"get_current_weather": getCurrentWeather,
		"get_stock_price":     getStockPrice,
	}

	// Create a new client with tools
	client := xai.NewClient(
		apiKey,
		string(xai.Grok3Beta),
		xai.WithTools(tools),
		xai.WithFunctionMap(functionMap),
	)

	// Create a message that will likely trigger tool usage
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": "What's the weather like in Shanghai? Also, what's the current price of AAPL stock?",
		},
	}

	// Make the API call with tool_choice set to auto
	resp, err := client.Invoke(
		messages,
		xai.WithToolChoice("auto"),
	)
	if err != nil {
		log.Fatalf("Error calling API: %v", err)
	}

	// Print the initial response
	fmt.Println("Initial response:")
	fmt.Println(resp.Choices[0].Message.Content)

	// Print tool calls if present
	if len(resp.Choices[0].Message.ToolCalls) > 0 {
		fmt.Println("\nTool calls:")
		for i, toolCall := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("  %d. %s: %s\n", i+1, toolCall.Function.Name, toolCall.Function.Arguments)
		}
	}

	// Print tool results if present
	if len(resp.Choices[0].Message.ToolResults) > 0 {
		fmt.Println("\nTool results:")
		for i, toolResult := range resp.Choices[0].Message.ToolResults {
			fmt.Printf("  %d. %s: %s\n", i+1, toolResult.Name, toolResult.Content)
		}
	}
}
