package security

import (
	"path/filepath"
	"testing"
)

func TestSandbox_Resolve_inside(t *testing.T) {
	dir := t.TempDir()
	s := NewSandbox(dir)
	resolved, err := s.Resolve("foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.HasPrefix(resolved, dir) {
		t.Errorf("resolved %s not under %s", resolved, dir)
	}
}

func TestSandbox_Resolve_outside(t *testing.T) {
	dir := t.TempDir()
	s := NewSandbox(dir)
	_, err := s.Resolve("../../../etc/passwd")
	if err != ErrPathOutside {
		t.Errorf("expected ErrPathOutside, got %v", err)
	}
}
