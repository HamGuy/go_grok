# 🧪 Examples

This directory contains simple examples showing how to use the xAI Grok SDK for Go. Each example demonstrates a specific feature or usage pattern.

## 🚀 Running the Examples

To run any example, navigate to its directory and use `go run main.go`. You'll need to set the `GROK_API_KEY` environment variable first.

```bash
# Set API Key environment variable
export GROK_API_KEY="your-xai-api-key"

# Run the example
cd examples/basic
go run main.go
```

You can also create a `.env` file and use [godotenv](https://github.com/joho/godotenv) to load environment variables:

```
# .env file example
GROK_API_KEY=your_api_key_here
```

```go
// Load .env file in your code
if err := godotenv.Load(); err != nil {
    log.Println("Warning: No .env file found")
}
```

## 📋 Available Examples

### 📝 Basic Usage ([basic/main.go](basic/main.go))

A simple example showing how to make a basic request to the Grok API and handle the response.

Features demonstrated:
- Client initialization
- Setting request options
- Making a basic API call
- Handling the response
- Displaying token usage

### 🔄 Streaming Responses ([streaming/main.go](streaming/main.go))

Shows how to stream responses from the Grok API in real-time.

Features demonstrated:
- Setting up streaming
- Processing streaming chunks
- Displaying a loading animation
- Handling the complete streamed response

### 🧰 Tool Calling ([tool_calling/main.go](tool_calling/main.go))

Demonstrates how to define and use tools that the Grok model can call.

Features demonstrated:
- Defining custom tools
- Implementing tool handler functions
- Setting up a client with tools
- Processing tool calls and results
- Making a request that triggers tool use

## 🔍 More Complex Examples

For more complex examples and demonstrations, check out the [`demo/`](../demo/) directory. 