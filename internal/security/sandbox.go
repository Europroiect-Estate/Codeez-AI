package security

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Blocked path prefixes (case-insensitive on Windows).
var blockedPrefixes = []string{
	"~/.ssh",
	"~/.gnupg",
	"/etc",
	"/System",
	"/Windows",
	"/Program Files",
	"/Program Files (x86)",
	"$HOME/.ssh",
	"$HOME/.gnupg",
	"$HOME/.config", // except codeez
}

var (
	ErrPathBlocked = errors.New("path is blocked by sandbox")
	ErrPathOutside = errors.New("path is outside allowed root")
)

// Sandbox validates paths against an allowed root and blocklist.
type Sandbox struct {
	AllowedRoot string
}

// NewSandbox returns a sandbox with the given allowed root (e.g. git root or cwd).
func NewSandbox(allowedRoot string) *Sandbox {
	if allowedRoot == "" {
		allowedRoot = "."
	}
	abs, _ := filepath.Abs(allowedRoot)
	return &Sandbox{AllowedRoot: abs}
}

// Resolve resolves path to an absolute path and checks it is inside AllowedRoot and not blocked.
func (s *Sandbox) Resolve(path string) (string, error) {
	if path == "" {
		return "", ErrPathOutside
	}
	// Expand ~ and $HOME
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}
	if strings.HasPrefix(path, "$HOME") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, strings.TrimPrefix(path, "$HOME"))
	}
	base := s.AllowedRoot
	if base == "." {
		base, _ = os.Getwd()
	}
	// If path is relative, resolve against allowed root
	if !filepath.IsAbs(path) && !strings.HasPrefix(path, "~") && !strings.HasPrefix(path, "$HOME") {
		path = filepath.Join(base, path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, abs)
	if err != nil {
		return "", ErrPathOutside
	}
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", ErrPathOutside
	}
	// Blocklist
	norm := filepath.ToSlash(abs)
	if runtime.GOOS == "windows" {
		norm = strings.ToLower(norm)
	}
	for _, prefix := range blockedPrefixes {
		p := prefix
		if strings.HasPrefix(p, "~") {
			home, _ := os.UserHomeDir()
			p = filepath.Join(home, strings.TrimPrefix(p, "~/"))
		}
		p = filepath.ToSlash(p)
		if runtime.GOOS == "windows" {
			p = strings.ToLower(p)
		}
		if strings.HasPrefix(norm, p) || norm == p {
			return "", ErrPathBlocked
		}
	}
	return abs, nil
}

// AllowRead returns true if the path is allowed for read (same as Resolve).
func (s *Sandbox) AllowRead(path string) (string, error) {
	return s.Resolve(path)
}

// AllowWrite returns true if the path is allowed for write (same as Resolve).
func (s *Sandbox) AllowWrite(path string) (string, error) {
	return s.Resolve(path)
}
