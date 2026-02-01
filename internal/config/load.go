package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "toml"
	globalDir  = ".config/codeez"
	projectDir = ".codeez"
)

// Load reads config from global then project; project overrides global.
// globalBase is e.g. os.UserHomeDir(); projectRoot is the repo or cwd (can be empty).
func Load(globalBase, projectRoot string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	// Defaults
	v.SetDefault("provider", ProviderOllama)
	v.SetDefault("model", "")
	v.SetDefault("palette", PaletteOriginal)
	v.SetDefault("paths.global_config_dir", "")
	v.SetDefault("paths.project_config_dir", "")

	globalPath := filepath.Join(globalBase, globalDir)
	v.AddConfigPath(globalPath)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading global config: %w", err)
		}
	}

	if projectRoot != "" {
		projectConfigFile := filepath.Join(projectRoot, projectDir, configName+"."+configType)
		v.SetConfigFile(projectConfigFile)
		_ = v.MergeInConfig() // project override; ignore not found
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if c.Paths.GlobalConfigDir == "" {
		c.Paths.GlobalConfigDir = filepath.Join(globalBase, globalDir)
	}
	if c.Paths.ProjectConfigDir == "" && projectRoot != "" {
		c.Paths.ProjectConfigDir = filepath.Join(projectRoot, projectDir)
	}

	// Provider-specific keys from viper (e.g. providers.ollama.api_key)
	if c.Providers == nil {
		c.Providers = make(map[string]ProviderConfig)
	}
	for _, name := range ValidProviders() {
		key := "providers." + name
		if v.IsSet(key + ".api_key") {
			if c.Providers[name].APIKey == "" {
				pc := c.Providers[name]
				pc.APIKey = v.GetString(key + ".api_key")
				pc.Model = v.GetString(key + ".model")
				pc.BaseURL = v.GetString(key + ".base_url")
				c.Providers[name] = pc
			}
		}
	}
	if c.APIKey == "" && c.Provider != "" {
		if p, ok := c.Providers[c.Provider]; ok && p.APIKey != "" {
			c.APIKey = p.APIKey
		}
	}
	if c.Model == "" && c.Provider != "" {
		if p, ok := c.Providers[c.Provider]; ok && p.Model != "" {
			c.Model = p.Model
		}
	}

	return &c, nil
}

// Save writes key-value to the given config file (global or project).
func Save(configPath, key, value string) error {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType(configType)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.Set(key, value)
	return v.WriteConfig()
}

// EnsureConfigDir creates the config directory and writes an empty config if missing.
func EnsureConfigDir(dir string) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	path := filepath.Join(dir, configName+"."+configType)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.WriteFile(path, []byte("# codeez config\n"), 0o600)
	}
	return nil
}
