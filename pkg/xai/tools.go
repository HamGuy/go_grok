package xai

import (
	"encoding/json"
	"fmt"
)

// ExecuteToolCall executes a single tool call and returns the result
func (c *Client) ExecuteToolCall(toolCall ToolCall) (*ToolResult, error) {
	if toolCall.Type != "function" {
		return nil, fmt.Errorf("unsupported tool call type: %s", toolCall.Type)
	}

	funcName := toolCall.Function.Name
	fn, ok := c.FunctionMap[funcName]
	if !ok {
		return nil, fmt.Errorf("function not found in function map: %s", funcName)
	}

	// Handle arguments based on their actual type
	var argsMap map[string]interface{}

	switch args := toolCall.Function.Arguments.(type) {
	case string:
		// If arguments is a JSON string, parse it
		if err := json.Unmarshal([]byte(args), &argsMap); err != nil {
			return nil, fmt.Errorf("failed to parse function arguments: %v", err)
		}
	case map[string]interface{}:
		// If arguments is already a map, use it directly
		argsMap = args
	default:
		return nil, fmt.Errorf("unexpected arguments type: %T", toolCall.Function.Arguments)
	}

	// Execute the function with the arguments
	result, err := fn(argsMap)
	if err != nil {
		return &ToolResult{
			ToolCallID: toolCall.ID,
			Role:       "tool",
			Name:       funcName,
			Content:    fmt.Sprintf("Error: %s", err.Error()),
		}, nil
	}

	return &ToolResult{
		ToolCallID: toolCall.ID,
		Role:       "tool",
		Name:       funcName,
		Content:    result,
	}, nil
}

// ValidateTools checks if all tools have valid names and schemas
func ValidateTools(tools []map[string]interface{}) error {
	for i, tool := range tools {
		name, ok := tool["name"].(string)
		if !ok || name == "" {
			return fmt.Errorf("tool at index %d must have a non-empty string name", i)
		}

		// Check if the tool has a valid schema
		if parameters, ok := tool["parameters"].(map[string]interface{}); ok {
			if typ, ok := parameters["type"].(string); !ok || typ != "object" {
				return fmt.Errorf("tool '%s' parameters must have type 'object'", name)
			}
		} else {
			return fmt.Errorf("tool '%s' must have valid parameters schema", name)
		}
	}

	return nil
}

// ValidateFunctionMap checks if all functions in the map have corresponding tools
func ValidateFunctionMap(tools []map[string]interface{}, functionMap map[string]func(map[string]interface{}) (string, error)) error {
	// Create a set of tool names
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		if name, ok := tool["name"].(string); ok {
			toolNames[name] = true
		}
	}

	// Check if all functions in the map have corresponding tools
	for funcName := range functionMap {
		if !toolNames[funcName] {
			return fmt.Errorf("function '%s' in function map has no corresponding tool definition", funcName)
		}
	}

	return nil
}
