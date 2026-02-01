package providers

import (
	"context"
	"fmt"
)

// OpenAIConfig configures the OpenAI provider.
type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// OpenAI implements Provider for OpenAI (stub/minimal for Phase C).
type OpenAI struct {
	apiKey string
	model  string
	base   string
}

// NewOpenAI returns an OpenAI provider.
func NewOpenAI(cfg OpenAIConfig) *OpenAI {
	base := cfg.BaseURL
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	return &OpenAI{apiKey: cfg.APIKey, model: cfg.Model, base: base}
}

// ChatStream streams from OpenAI (stub: returns error until implemented).
func (o *OpenAI) ChatStream(ctx context.Context, messages []Message, toolsSchema interface{}) (Stream, error) {
	if o.apiKey == "" {
		return nil, fmt.Errorf("openai: api_key not set")
	}
	// TODO: implement OpenAI chat completions streaming
	return nil, fmt.Errorf("openai: streaming not implemented yet")
}
