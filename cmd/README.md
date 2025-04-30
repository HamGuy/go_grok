# 🖥️ Command-Line Tools

This directory contains command-line applications built with the xAI Grok SDK for Go. These tools demonstrate how to use the SDK in practical applications.

## 🚀 Server Example

### 📡 API Server ([server/](server/))

A simple HTTP server that acts as a proxy to the xAI API, allowing you to:

- Expose a local API endpoint for Grok interactions
- Add middleware for authentication, logging, rate limiting, etc.
- Cache responses or implement custom behavior
- Create a unified interface for internal applications

### 💡 Environment Variable Configuration

The server is configured using environment variables:

```bash
# Required: Grok API Key
export GROK_API_KEY="your-grok-api-key"

# Optional: Server port (default: 8080)
export PORT="8080"

# Run the server
cd cmd/server
go run main.go
```

You can also create a `.env` file:

```
# .env file example
GROK_API_KEY=your_grok_api_key_here
PORT=8080
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

### 🔍 Usage

After starting the server, you can send requests to `http://localhost:8080/chat/completions`:

```bash
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [{"role": "user", "content": "Hello, how are you?"}],
    "temperature": 0.7,
    "max_tokens": 500
  }'
```

Streaming response:

```bash
curl -X POST http://localhost:8080/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [{"role": "user", "content": "Write a short poem"}],
    "stream": true
  }'
```

## 🔗 Additional Resources

- For simple usage examples, check out the [`examples/`](../examples/) directory
- For more complex demonstrations, see the [`demo/`](../demo/) directory 