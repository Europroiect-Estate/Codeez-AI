package providers

import (
	"context"
	"fmt"
)

// AnthropicConfig configures the Anthropic provider.
type AnthropicConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// Anthropic implements Provider for Anthropic (stub/minimal for Phase C).
type Anthropic struct {
	apiKey string
	model  string
	base   string
}

// NewAnthropic returns an Anthropic provider.
func NewAnthropic(cfg AnthropicConfig) *Anthropic {
	base := cfg.BaseURL
	if base == "" {
		base = "https://api.anthropic.com"
	}
	return &Anthropic{apiKey: cfg.APIKey, model: cfg.Model, base: base}
}

// ChatStream streams from Anthropic (stub: returns error until implemented).
func (a *Anthropic) ChatStream(ctx context.Context, messages []Message, toolsSchema interface{}) (Stream, error) {
	if a.apiKey == "" {
		return nil, fmt.Errorf("anthropic: api_key not set")
	}
	// TODO: implement Anthropic messages streaming
	return nil, fmt.Errorf("anthropic: streaming not implemented yet")
}
