package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// RepoMapTool returns key files, modules, languages, build hints (from .codeez/repo_map.json or best-effort scan).
type RepoMapTool struct {
	Cwd string
}

func (t *RepoMapTool) Name() string { return "repo.map" }
func (t *RepoMapTool) Description() string {
	return "Return key files, modules, languages, build hints"
}
func (t *RepoMapTool) Parameters() string { return `{}` }

func (t *RepoMapTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	cwd := t.Cwd
	if cwd == "" {
		cwd = "."
	}
	// Try .codeez/repo_map.json first (from codeez index)
	path := filepath.Join(cwd, ".codeez", "repo_map.json")
	data, err := os.ReadFile(path)
	if err == nil {
		var m map[string]interface{}
		if json.Unmarshal(data, &m) == nil {
			b, _ := json.MarshalIndent(m, "", "  ")
			return string(b), nil
		}
	}
	// Best-effort: key files by extension
	keyFiles := []string{}
	_ = filepath.Walk(cwd, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(cwd, p)
		if strings.Contains(rel, ".git") {
			return filepath.SkipDir
		}
		switch filepath.Ext(p) {
		case ".go":
			if strings.HasSuffix(rel, "go.mod") || strings.HasSuffix(rel, "main.go") {
				keyFiles = append(keyFiles, rel)
			}
		case ".json":
			if strings.Contains(rel, "package.json") {
				keyFiles = append(keyFiles, rel)
			}
		}
		return nil
	})
	out := map[string]interface{}{
		"key_files": keyFiles,
		"languages": []string{"go"},
		"build":     "go build ./... or npm install",
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b), nil
}
