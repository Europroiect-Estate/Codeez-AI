package providers

import (
	"github.com/Europroiect-Estate/Codeez-AI/internal/config"
)

// Factory returns a Provider by name from config.
type Factory struct {
	Config *config.Config
}

// NewFactory returns a factory that uses the given config.
func NewFactory(cfg *config.Config) *Factory {
	return &Factory{Config: cfg}
}

// Get returns the active provider (from config.Provider) or nil if not configured.
func (f *Factory) Get() (Provider, error) {
	name := f.Config.Provider
	if name == "" {
		name = "ollama"
	}
	switch name {
	case "ollama":
		base := ""
		model := f.Config.Model
		if p, ok := f.Config.Providers["ollama"]; ok {
			if p.BaseURL != "" {
				base = p.BaseURL
			}
			if p.Model != "" {
				model = p.Model
			}
		}
		if base == "" {
			base = "http://localhost:11434"
		}
		return NewOllama(OllamaConfig{BaseURL: base, Model: model}), nil
	case "openai":
		apiKey := f.Config.APIKey
		model := f.Config.Model
		if p, ok := f.Config.Providers["openai"]; ok {
			if p.APIKey != "" {
				apiKey = p.APIKey
			}
			if p.Model != "" {
				model = p.Model
			}
		}
		return NewOpenAI(OpenAIConfig{APIKey: apiKey, Model: model, BaseURL: ""}), nil
	case "anthropic":
		apiKey := f.Config.APIKey
		model := f.Config.Model
		if p, ok := f.Config.Providers["anthropic"]; ok {
			if p.APIKey != "" {
				apiKey = p.APIKey
			}
			if p.Model != "" {
				model = p.Model
			}
		}
		return NewAnthropic(AnthropicConfig{APIKey: apiKey, Model: model, BaseURL: ""}), nil
	default:
		return NewOllama(OllamaConfig{BaseURL: "http://localhost:11434", Model: f.Config.Model}), nil
	}
}
