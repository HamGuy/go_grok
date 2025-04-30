// Package utils provides utility functions for the SDK
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadEnvFile loads the .env file from the project root or current directory
func LoadEnvFile() {
	// Try to load from current directory
	loadEnvFromFile(".env")

	// Try to load from project root
	rootDir := findProjectRoot()
	if rootDir != "" {
		loadEnvFromFile(filepath.Join(rootDir, ".env"))
	}
}

// GetAPIKey gets the API key, prioritizes environment variable, then tries .env file
func GetAPIKey() string {
	// Get from environment variable first
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	// If environment variable is not set, try to load from .env file
	LoadEnvFile()

	// Try to get from environment variable again
	apiKey = os.Getenv("GROK_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	// If still not set, try to read from file
	rootDir := findProjectRoot()
	if rootDir != "" {
		apiKey = extractAPIKeyFromFile(filepath.Join(rootDir, ".env"))
	}

	// If still not found, try to find in current directory
	if apiKey == "" {
		apiKey = extractAPIKeyFromFile(".env")
	}

	return apiKey
}

// GetMaskedAPIKey gets a masked version of the API key
func GetMaskedAPIKey(key string) string {
	if key == "" {
		return "****"
	}
	if len(key) <= 8 {
		return "****"
	}
	// Only show the first 4 and last 4 characters, replace the middle with ****
	return key[:4] + "****" + key[len(key)-4:]
}

// Extract API key from file
func extractAPIKeyFromFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "GROK_API_KEY=") {
			return strings.TrimPrefix(trimmedLine, "GROK_API_KEY=")
		}
	}

	return ""
}

// Load environment variables from .env file
func loadEnvFromFile(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Skip empty lines and comment lines
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(trimmedLine, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set environment variable
		if key != "" && value != "" {
			os.Setenv(key, value)
		}
	}
}

// Find project root directory
func findProjectRoot() string {
	// Start from current directory and search upwards for go.mod file
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		// Check if current directory has go.mod file
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Go up one level
		parent := filepath.Dir(dir)
		if parent == dir {
			// Already at root, go.mod not found
			return ""
		}
		dir = parent
	}
}
