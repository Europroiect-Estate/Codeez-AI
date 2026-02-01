package app

import (
	"os"
	"path/filepath"

	"github.com/Europroiect-Estate/Codeez-AI/internal/config"
)

// App holds paths and config; no heavy deps.
type App struct {
	Config     *config.Config
	Cwd        string
	RepoRoot   string
	GlobalBase string
}

// New loads config and sets cwd/repo root. repoRoot can be empty.
// Project config is merged when .codeez exists in cwd or in repoRoot.
func New(globalBase, cwd, repoRoot string) (*App, error) {
	if globalBase == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		globalBase = home
	}
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	projectRootForConfig := repoRoot
	if projectRootForConfig == "" {
		if _, err := os.Stat(filepath.Join(cwd, ".codeez")); err == nil {
			projectRootForConfig = cwd
		}
	}

	cfg, err := config.Load(globalBase, projectRootForConfig)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:     cfg,
		Cwd:        cwd,
		RepoRoot:   repoRoot,
		GlobalBase: globalBase,
	}, nil
}

// GlobalConfigDir returns the global config directory.
func (a *App) GlobalConfigDir() string {
	if a.Config != nil && a.Config.Paths.GlobalConfigDir != "" {
		return a.Config.Paths.GlobalConfigDir
	}
	return filepath.Join(a.GlobalBase, ".config", "codeez")
}

// ProjectConfigDir returns the project config directory (empty if not in a project).
// Prefer .codeez in Cwd if it exists; else use repo root.
func (a *App) ProjectConfigDir() string {
	if a.Config != nil && a.Config.Paths.ProjectConfigDir != "" {
		return a.Config.Paths.ProjectConfigDir
	}
	cwdDot := filepath.Join(a.Cwd, ".codeez")
	if _, err := os.Stat(cwdDot); err == nil {
		return cwdDot
	}
	if a.RepoRoot != "" {
		return filepath.Join(a.RepoRoot, ".codeez")
	}
	return ""
}
