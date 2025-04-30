package xai

// InvokeOption is a function that configures a ChatCompletionRequest
type InvokeOption func(*ChatCompletionRequest)

// WithTemperature sets the temperature parameter
func WithTemperature(temp float64) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.Temperature = &temp
	}
}

// WithMaxTokens sets the max_tokens parameter
func WithMaxTokens(tokens int) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.MaxTokens = &tokens
	}
}

// WithStream enables or disables streaming
func WithStream(stream bool) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.Stream = &stream
	}
}

// WithToolChoice sets the tool_choice parameter
func WithToolChoice(choice interface{}) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.ToolChoice = choice
	}
}

// WithFrequencyPenalty sets the frequency_penalty parameter
func WithFrequencyPenalty(penalty float64) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.FrequencyPenalty = &penalty
	}
}

// WithPresencePenalty sets the presence_penalty parameter
func WithPresencePenalty(penalty float64) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.PresencePenalty = &penalty
	}
}

// WithTopP sets the top_p parameter
func WithTopP(topP float64) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.TopP = &topP
	}
}

// WithSeed sets the seed parameter
func WithSeed(seed int) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.Seed = &seed
	}
}

// WithN sets the n parameter
func WithN(n int) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.N = &n
	}
}

// WithStop sets the stop parameter
func WithStop(stop []string) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.Stop = stop
	}
}

// WithLogProbs enables or disables log probabilities
func WithLogProbs(logprobs bool) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.LogProbs = &logprobs
	}
}

// WithTopLogProbs sets the top_logprobs parameter
func WithTopLogProbs(topLogProbs int) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.TopLogProbs = &topLogProbs
	}
}

// WithResponseFormat sets the response_format parameter
func WithResponseFormat(format interface{}) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.ResponseFormat = format
	}
}

// WithUser sets the user parameter
func WithUser(user string) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.User = &user
	}
}

// WithLogitBias sets the logit_bias parameter
func WithLogitBias(logitBias map[string]interface{}) InvokeOption {
	return func(req *ChatCompletionRequest) {
		req.LogitBias = logitBias
	}
}
