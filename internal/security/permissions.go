package security

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// PermissionsFile holds "always allow for this repo" rules keyed by repo root.
type PermissionsFile struct {
	RepoRoot string            `toml:"repo_root"`
	Rules    []PermissionRule  `toml:"rules"`
}

// PermissionRule is a single allow rule (tool + optional pattern + scope).
type PermissionRule struct {
	ToolName string `toml:"tool_name"`
	Pattern  string `toml:"pattern,omitempty"`
	Scope    string `toml:"scope"` // e.g. "repo"
}

const permissionsFileName = "permissions.toml"

// LoadPermissions reads .codeez/permissions.toml from the given dir.
func LoadPermissions(dotCodeezDir string) (*PermissionsFile, error) {
	if dotCodeezDir == "" {
		return &PermissionsFile{}, nil
	}
	path := filepath.Join(dotCodeezDir, permissionsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &PermissionsFile{}, nil
		}
		return nil, err
	}
	var f PermissionsFile
	if _, err := toml.Decode(string(data), &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// SavePermissions writes .codeez/permissions.toml.
func SavePermissions(dotCodeezDir string, f *PermissionsFile) error {
	if dotCodeezDir == "" {
		return nil
	}
	if err := os.MkdirAll(dotCodeezDir, 0o750); err != nil {
		return err
	}
	path := filepath.Join(dotCodeezDir, permissionsFileName)
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(f); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o600)
}

// AddRule appends a rule and saves. It loads, appends, and saves.
func AddRule(dotCodeezDir, repoRoot, toolName, pattern, scope string) error {
	f, err := LoadPermissions(dotCodeezDir)
	if err != nil {
		return err
	}
	f.RepoRoot = repoRoot
	f.Rules = append(f.Rules, PermissionRule{ToolName: toolName, Pattern: pattern, Scope: scope})
	return SavePermissions(dotCodeezDir, f)
}

// HasRule returns true if there is an "always allow" rule for the given tool (and optional pattern).
func (f *PermissionsFile) HasRule(toolName, pattern string) bool {
	for _, r := range f.Rules {
		if r.ToolName != toolName {
			continue
		}
		if pattern == "" || r.Pattern == "" || r.Pattern == pattern {
			return true
		}
	}
	return false
}
