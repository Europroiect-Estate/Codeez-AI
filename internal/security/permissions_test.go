package security

import (
	"path/filepath"
	"testing"
)

func TestPermissions_roundtrip(t *testing.T) {
	dir := t.TempDir()
	dotCodeez := filepath.Join(dir, ".codeez")
	if err := AddRule(dotCodeez, "/repo", "fs.read", "", "repo"); err != nil {
		t.Fatal(err)
	}
	f, err := LoadPermissions(dotCodeez)
	if err != nil {
		t.Fatal(err)
	}
	if f.RepoRoot != "/repo" || len(f.Rules) != 1 || f.Rules[0].ToolName != "fs.read" {
		t.Errorf("LoadPermissions: got %+v", f)
	}
	if !f.HasRule("fs.read", "") {
		t.Error("HasRule fs.read: expected true")
	}
}
