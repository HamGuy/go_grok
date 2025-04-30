# 🎮 Advanced Demos

This directory contains more complex and complete demonstrations of the xAI Grok SDK for Go. These examples showcase real-world use cases and advanced features beyond the simple examples.

## 🚀 Running the Demos

To run any demo, navigate to its directory and use `go run main.go`. You'll need to set the `GROK_API_KEY` environment variable first.

```bash
# Set API Key environment variable
export GROK_API_KEY="your-xai-api-key"

# Run the demo
cd demo/streaming
go run main.go
```

You can also use [godotenv](https://github.com/joho/godotenv) to load environment variables from a `.env` file:

```bash
# Install godotenv
go get github.com/joho/godotenv
```

```
# .env file example
GROK_API_KEY=your_api_key_here
```

## 📋 Available Demos

### 🔄 Advanced Streaming ([streaming/main.go](streaming/main.go))

A more advanced demonstration of streaming capabilities with the Grok API, featuring:

- Real-time output display with cursor control
- Loading spinner animation while waiting for the first response
- Efficient management of streaming goroutines
- Complete response collection alongside streaming display

This demo shows how to create a polished user experience with streaming responses, suitable for CLI applications or more complex integration scenarios.

## 🔗 Additional Resources

- For simpler examples, check out the [`examples/`](../examples/) directory
- For command-line applications, see the [`cmd/`](../cmd/) directory 