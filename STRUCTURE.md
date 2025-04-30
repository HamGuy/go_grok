# 📁 Project Structure

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
│   ├── utils/              # Utility functions
│   │   └── env.go          # Environment and configuration utilities
│   └── xai/                # Main SDK package
│       ├── client.go       # Main client implementation
│       ├── client_test.go  # Tests for client functionality
│       ├── models.go       # Data models and structures
│       ├── options.go      # Configuration options
│       ├── streaming.go    # Streaming response handling
│       └── tools.go        # Tool calling functionality
├── main.go                 # Root example file
├── go.mod                  # Go module definition
└── README.md               # Project documentation
```

## 📂 Directory Overview

### 📁 cmd/
The `cmd/` directory contains command-line applications that use the SDK. These serve as both examples and useful tools:
- `server/`: A simple API server that acts as a proxy to the xAI API

### 📁 demo/
The `demo/` directory hosts more complex example applications:
- `streaming/`: Advanced demonstrations of streaming capabilities

### 📁 examples/
The `examples/` directory contains simple, focused examples:
- `basic/`: Illustrates basic SDK usage for simple chat completions
- `streaming/`: Shows how to use streaming responses
- `tool_calling/`: Demonstrates function calling capabilities

### 📁 pkg/
The `pkg/` directory contains the core SDK code:
- `utils/`: Utility functions for the SDK
  - `env.go`: Utilities for handling environment variables and configuration
- `xai/`: The main package implementing the SDK functionality
  - `client.go`: Main client implementation for making API calls
  - `client_test.go`: Unit tests for client functionality
  - `models.go`: Data structures for requests and responses
  - `options.go`: Configuration options for API requests
  - `streaming.go`: Implementation of streaming functionality
  - `tools.go`: Implementation of tool calling functionality

### 📄 main.go
Root example file showing basic usage of the SDK with streaming capabilities. 