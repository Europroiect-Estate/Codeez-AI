package index

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RepoMap holds key files, modules, languages, build hints.
type RepoMap struct {
	KeyFiles  []string          `json:"key_files"`
	Modules   []string          `json:"modules,omitempty"`
	Languages []string          `json:"languages"`
	Build     string            `json:"build,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

const repoMapFile = "repo_map.json"

// Save writes the repo map to .codeez/repo_map.json.
func (r *RepoMap) Save(dotCodeezDir string) error {
	if dotCodeezDir == "" {
		return nil
	}
	if err := os.MkdirAll(dotCodeezDir, 0o750); err != nil {
		return err
	}
	path := filepath.Join(dotCodeezDir, repoMapFile)
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

// Load reads the repo map from .codeez/repo_map.json.
func Load(dotCodeezDir string) (*RepoMap, error) {
	if dotCodeezDir == "" {
		return &RepoMap{}, nil
	}
	path := filepath.Join(dotCodeezDir, repoMapFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &RepoMap{}, nil
		}
		return nil, err
	}
	var r RepoMap
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
