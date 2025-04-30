# go_grok -  🤖 xAI Grok SDK for Go

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/xAI-Grok-6236FF?style=for-the-badge&logo=x&logoColor=white" alt="xAI Grok">
</p>

## 📝 Overview

go-grok is a lightweight, feature-complete Go SDK for interacting with xAI's Grok API. This SDK provides a simple, intuitive interface to leverage the capabilities of Grok models in your Go applications.

### ✨ Supported Models

- **Grok 2 Series** 🧠
  - `grok-2-1212`
  - `grok-2-vision-1212`
  - `grok-2-image-1212`

- **Grok 3 Series** 🚀
  - `grok-3-beta`
  - `grok-3-fast-beta`
  - `grok-3-mini-beta`
  - `grok-3-mini-fast-beta`

## 🛠️ Features

- ✅ **Basic Chat Completions** - Send messages and receive AI-generated responses
- ✅ **Streaming Support** - Stream responses in real-time
- ✅ **Tool Calling** - Define custom tools that Grok can use
- ✅ **Function Execution** - Automatically execute functions when Grok calls them
- ✅ **Customizable Parameters** - Full control over temperature, max tokens, etc.

## 📦 Installation

```bash
go get github.com/hamguy/go_grok
```

## 🚀 Quick Start

```go
package main

import (
	"fmt"
	"github.com/hamguy/go_grok/pkg/xai"
)

func main() {
	// Initialize client with your API key
	client := xai.NewClient(
		"YOUR_XAI_API_KEY",
		string(xai.Grok3Beta),
	)

	// Create a simple chat message
	messages := []map[string]interface{}{
		{"role": "user", "content": "Hello! What can you tell me about yourself?"},
	}

	// Get a response
	resp, err := client.Invoke(messages)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the response
	fmt.Println(resp.Choices[0].Message.Content)
}
```

## 📊 Advanced Usage

### 🔄 Streaming Responses

```go
// Enable streaming
streamChan, err := client.InvokeStream(
	[]map[string]interface{}{
		{"role": "user", "content": "Write a short poem about AI"},
	},
)

// Process the stream
for chunk := range streamChan {
	if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
		fmt.Print(chunk.Choices[0].Delta.Content)
	}
}
```

### 🧰 Using Tools

```go
// Define a weather tool
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

// Implement the function
getWeather := func(args map[string]interface{}) (string, error) {
	location, _ := args["location"].(string)
	return "The weather in " + location + " is sunny.", nil
}

// Create a function map
functionMap := map[string]func(map[string]interface{}) (string, error){
	"get_weather": getWeather,
}

// Create a client with tools
client := xai.NewClient(
	"YOUR_XAI_API_KEY",
	string(xai.Grok3Beta),
	xai.WithTools(tools),
	xai.WithFunctionMap(functionMap),
)

// Send a message that might trigger tool use
resp, err := client.Invoke(
	[]map[string]interface{}{
		{"role": "user", "content": "What's the weather in San Francisco?"},
	},
	xai.WithToolChoice("auto"),
)
```

### ⚙️ Configuring Parameters

```go
resp, err := client.Invoke(
	messages,
	xai.WithTemperature(0.7),
	xai.WithMaxTokens(1000),
	xai.WithTopP(0.95),
)
```

## 📚 API Reference

For detailed API reference, check our [API Documentation](https://github.com/hamguy/go_grok/wiki).

## 🧪 Testing

```bash
go test ./pkg/xai/
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- [xAI](https://x.ai) for creating the Grok AI models