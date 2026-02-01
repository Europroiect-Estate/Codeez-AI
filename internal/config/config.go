package config

// Config holds merged global and project configuration.
type Config struct {
	Provider string            `mapstructure:"provider"`
	Model    string            `mapstructure:"model"`
	APIKey   string            `mapstructure:"api_key"`
	Palette  string            `mapstructure:"palette"`
	Paths    PathsConfig       `mapstructure:"paths"`
	Providers map[string]ProviderConfig `mapstructure:"providers"`
}

// PathsConfig holds path-related settings.
type PathsConfig struct {
	GlobalConfigDir string `mapstructure:"global_config_dir"`
	ProjectConfigDir string `mapstructure:"project_config_dir"`
}

// ProviderConfig holds provider-specific settings (e.g. base_url, api_key).
type ProviderConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

const (
	PaletteOriginal  = "original"
	PaletteCorporate = "corporate"
	PaletteCyber     = "cyber"
)

const (
	ProviderOllama    = "ollama"
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
)

// ValidPalettes returns allowed palette values.
func ValidPalettes() []string {
	return []string{PaletteOriginal, PaletteCorporate, PaletteCyber}
}

// ValidProviders returns supported provider names.
func ValidProviders() []string {
	return []string{ProviderOllama, ProviderOpenAI, ProviderAnthropic}
}
