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
- ✅ **Server Integration** - Run as an API server with built-in streaming support

## 📦 Installation

```bash
go get github.com/hamguy/go_grok
```

## 🔑 API Key Configuration

This SDK requires an xAI Grok API key to function. You can provide it in several ways:

### Environment Variables

The recommended way is to use the `GROK_API_KEY` environment variable:

```bash
# Set API Key environment variable
export GROK_API_KEY="your-grok-api-key"
```

### Using .env Files

You can also use a `.env` file with [godotenv](https://github.com/joho/godotenv):

```bash
# Install godotenv
go get github.com/joho/godotenv
```

Create a `.env` file in your project:
```
# .env file example
GROK_API_KEY=your_api_key_here
```

Then load it in your code:
```go
import "github.com/joho/godotenv"

func init() {
    // Load .env file (if it exists)
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}
```

## 🚀 Quick Start

```go
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
		log.Fatal("GROK_API_KEY environment variable not set")
	}

	// Initialize client with your API key
	client := xai.NewClient(
		apiKey,
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
	apiKey,
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

### 🖥️ Server Integration

You can run this SDK as an API server that proxies requests to the xAI API:

```go
package main

import (
	"log"
	"os"
	"github.com/hamguy/go_grok/pkg/xai"
	"net/http"
	// Other imports...
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set")
	}
	
	// Create global client
	client := xai.NewClient(apiKey, string(xai.Grok3Beta))
	
	// Set up HTTP handlers
	http.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		// Handle chat completions endpoint
		// See cmd/server/main.go for complete implementation
	})
	
	// Start server
	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
```

For a complete server example, see [cmd/server/main.go](cmd/server/main.go).

## 📁 Project Structure

```
go_grok/  
├── cmd/                    # Command-line tools  
│   └── server/             # API server example
├── demo/                   # More complex demonstrations  
│   └── streaming/          # Streaming examples  
│       └── main.go         # Advanced streaming output example
├── examples/               # Simple usage examples  
│   ├── basic/main.go       # Basic usage example  
│   ├── streaming/main.go   # Simple streaming example  
│   └── tool_calling/main.go # Tool calling example  
├── pkg/                    # Core SDK code  
│   └── xai/                # Main package implementing the SDK
├── go.mod                  # Go module definition
└── README.md               # Project documentation
```

## 📚 Examples and Demos

### 📋 Simple Examples

The [`examples/`](examples/) directory contains simple, focused examples:

- **Basic Usage**: Simple chat completion requests
- **Streaming**: Real-time streaming of responses
- **Tool Calling**: Using function calling capabilities

### 🎮 Advanced Demos

The [`demo/`](demo/) directory contains more complex demonstrations:

- **Advanced Streaming**: Real-time output with loading animations and cursor control

### 🖥️ Command-Line Tools

The [`cmd/`](cmd/) directory contains command-line applications:

- **API Server**: HTTP server acting as a proxy to the xAI API

## 🧪 Testing

```bash
go test ./pkg/xai/
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- [xAI](https://x.ai) for creating the Grok AI models