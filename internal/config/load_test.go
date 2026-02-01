package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_GlobalOnly(t *testing.T) {
	dir := t.TempDir()
	globalBase := filepath.Join(dir, "home")
	globalConfig := filepath.Join(globalBase, ".config", "codeez")
	if err := os.MkdirAll(globalConfig, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(globalConfig, "config.toml"), []byte(`
provider = "ollama"
palette = "original"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(globalBase, "")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Provider != "ollama" {
		t.Errorf("provider: got %q", cfg.Provider)
	}
	if cfg.Palette != "original" {
		t.Errorf("palette: got %q", cfg.Palette)
	}
}

func TestLoad_ProjectOverride(t *testing.T) {
	dir := t.TempDir()
	globalBase := filepath.Join(dir, "home")
	globalConfig := filepath.Join(globalBase, ".config", "codeez")
	projectRoot := filepath.Join(dir, "proj")
	projectConfig := filepath.Join(projectRoot, ".codeez")
	for _, d := range []string{globalConfig, projectConfig} {
		if err := os.MkdirAll(d, 0o750); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(globalConfig, "config.toml"), []byte(`
provider = "ollama"
palette = "original"
`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectConfig, "config.toml"), []byte(`
palette = "corporate"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(globalBase, projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Provider != "ollama" {
		t.Errorf("provider: got %q", cfg.Provider)
	}
	if cfg.Palette != "corporate" {
		t.Errorf("palette: want corporate, got %q", cfg.Palette)
	}
}
