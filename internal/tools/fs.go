package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Europroiect-Estate/Codeez-AI/internal/security"
)

type FSTool struct {
	Sandbox *security.Sandbox
}

func (t *FSTool) Name() string        { return "fs.read" }
func (t *FSTool) Description() string { return "Read file content" }
func (t *FSTool) Parameters() string  { return `{"path": "string"}` }

func (t *FSTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path required")
	}
	resolved, err := t.Sandbox.AllowRead(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type FSWriteTool struct {
	Sandbox *security.Sandbox
}

func (t *FSWriteTool) Name() string        { return "fs.write" }
func (t *FSWriteTool) Description() string { return "Write file content" }
func (t *FSWriteTool) Parameters() string  { return `{"path": "string", "content": "string"}` }

func (t *FSWriteTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	path, _ := args["path"].(string)
	content, _ := args["content"].(string)
	if path == "" {
		return "", fmt.Errorf("path required")
	}
	resolved, err := t.Sandbox.AllowWrite(path)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(resolved), 0o750); err != nil {
		return "", err
	}
	if err := os.WriteFile(resolved, []byte(content), 0o644); err != nil {
		return "", err
	}
	return "ok", nil
}

type FSListTool struct {
	Sandbox *security.Sandbox
}

func (t *FSListTool) Name() string        { return "fs.list" }
func (t *FSListTool) Description() string { return "List directory entries" }
func (t *FSListTool) Parameters() string  { return `{"dir": "string"}` }

func (t *FSListTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	dir, _ := args["dir"].(string)
	if dir == "" {
		dir = "."
	}
	resolved, err := t.Sandbox.AllowRead(dir)
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(resolved)
	if err != nil {
		return "", err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return strings.Join(names, "\n"), nil
}

type FSSearchTool struct {
	Sandbox *security.Sandbox
}

func (t *FSSearchTool) Name() string        { return "fs.search" }
func (t *FSSearchTool) Description() string { return "Search file content (simple substring)" }
func (t *FSSearchTool) Parameters() string  { return `{"dir": "string", "query": "string"}` }

func (t *FSSearchTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	dir, _ := args["dir"].(string)
	query, _ := args["query"].(string)
	if dir == "" {
		dir = "."
	}
	if query == "" {
		return "", fmt.Errorf("query required")
	}
	resolved, err := t.Sandbox.AllowRead(dir)
	if err != nil {
		return "", err
	}
	var matches []string
	err = filepath.Walk(resolved, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), query) {
			rel, _ := filepath.Rel(resolved, p)
			matches = append(matches, rel)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return strings.Join(matches, "\n"), nil
}

type FSStatTool struct {
	Sandbox *security.Sandbox
}

func (t *FSStatTool) Name() string        { return "fs.stat" }
func (t *FSStatTool) Description() string { return "Stat a path" }
func (t *FSStatTool) Parameters() string  { return `{"path": "string"}` }

func (t *FSStatTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path required")
	}
	resolved, err := t.Sandbox.AllowRead(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("mode=%v size=%d dir=%v", info.Mode(), info.Size(), info.IsDir()), nil
}
