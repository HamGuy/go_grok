// Package xai provides a client for the xAI API
package xai

import (
	"fmt"
)

// ModelType represents supported xAI models
type ModelType string

const (
	// Grok 2 系列
	Grok212     ModelType = "grok-2-1212"
	Grok2Vision ModelType = "grok-2-vision-1212"
	Grok2Image  ModelType = "grok-2-image-1212"

	// Grok 3 系列
	Grok3Beta         ModelType = "grok-3-beta"
	Grok3FastBeta     ModelType = "grok-3-fast-beta"
	Grok3MiniBeta     ModelType = "grok-3-mini-beta"
	Grok3MiniFastBeta ModelType = "grok-3-mini-fast-beta"
)

// ChatCompletionRequest represents a request to the chat completions API
type ChatCompletionRequest struct {
	Messages         []map[string]interface{} `json:"messages"`
	Model            string                   `json:"model"`
	FrequencyPenalty *float64                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{}   `json:"logit_bias,omitempty"`
	LogProbs         *bool                    `json:"logprobs,omitempty"`
	MaxTokens        *int                     `json:"max_tokens,omitempty"`
	N                *int                     `json:"n,omitempty"`
	PresencePenalty  *float64                 `json:"presence_penalty,omitempty"`
	ResponseFormat   interface{}              `json:"response_format,omitempty"`
	Seed             *int                     `json:"seed,omitempty"`
	Stop             []string                 `json:"stop,omitempty"`
	Stream           *bool                    `json:"stream,omitempty"`
	StreamOptions    interface{}              `json:"stream_options,omitempty"`
	Temperature      *float64                 `json:"temperature,omitempty"`
	Tools            []map[string]interface{} `json:"tools,omitempty"`
	ToolChoice       interface{}              `json:"tool_choice,omitempty"`
	TopLogProbs      *int                     `json:"top_logprobs,omitempty"`
	TopP             *float64                 `json:"top_p,omitempty"`
	User             *string                  `json:"user,omitempty"`
}

// Validate checks and normalizes the request parameters
func (r *ChatCompletionRequest) Validate() error {
	if len(r.Messages) == 0 {
		return fmt.Errorf("messages is required and cannot be empty")
	}
	if r.Model == "" {
		return fmt.Errorf("model is required and cannot be empty")
	}

	// Initialize defaults
	if r.LogitBias == nil {
		r.LogitBias = make(map[string]interface{})
	}
	if r.Stop == nil {
		r.Stop = []string{}
	}
	if r.Tools == nil {
		r.Tools = []map[string]interface{}{}
	}

	// Clamp values to valid ranges
	if r.FrequencyPenalty != nil {
		*r.FrequencyPenalty = clampFloat(*r.FrequencyPenalty, -2.0, 2.0)
	}
	if r.PresencePenalty != nil {
		*r.PresencePenalty = clampFloat(*r.PresencePenalty, -2.0, 2.0)
	}
	if r.Temperature != nil {
		*r.Temperature = clampFloat(*r.Temperature, 0.0, 2.0)
	}
	if r.TopP != nil {
		*r.TopP = clampFloat(*r.TopP, 0.0, 1.0)
	}
	if r.TopLogProbs != nil {
		*r.TopLogProbs = clampInt(*r.TopLogProbs, 0, 20)
	}

	// Ensure stop sequences don't exceed limit
	if len(r.Stop) > 4 {
		r.Stop = r.Stop[:4]
	}

	return nil
}

// Helper functions for clamping values
func clampFloat(value float64, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func clampInt(value int, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Response models
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int64    `json:"created,omitempty"`
	Model             string   `json:"model,omitempty"`
	Object            string   `json:"object,omitempty"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Usage             *Usage   `json:"usage,omitempty"`
}

type Choice struct {
	Index        int                    `json:"index"`
	Message      *Message               `json:"message,omitempty"`
	Delta        *Message               `json:"delta,omitempty"`
	FinishReason string                 `json:"finish_reason,omitempty"`
	LogProbs     map[string]interface{} `json:"logprobs,omitempty"`
}

type Message struct {
	Role        string       `json:"role"`
	Content     string       `json:"content"`
	ToolCalls   []ToolCall   `json:"tool_calls,omitempty"`
	ToolResults []ToolResult `json:"tool_results,omitempty"`
	Refusal     interface{}  `json:"refusal,omitempty"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Function Function `json:"function"`
	Type     string   `json:"type"`
}

type Function struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments"` // Changed from map[string]interface{} to interface{}
}

type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Role       string `json:"role"`
	Name       string `json:"name"`
	Content    string `json:"content"`
}

type Usage struct {
	PromptTokens        int                    `json:"prompt_tokens"`
	CompletionTokens    int                    `json:"completion_tokens"`
	TotalTokens         int                    `json:"total_tokens"`
	PromptTokensDetails map[string]interface{} `json:"prompt_tokens_details,omitempty"`
}
