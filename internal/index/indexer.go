package index

import (
	"os"
	"path/filepath"
	"strings"
)

// Indexer scans the repo and builds a RepoMap (best-effort without ripgrep/tree-sitter for MVP).
type Indexer struct {
	Root string
}

// Index walks the repo (respecting .gitignore hints) and builds the map.
func (i *Indexer) Index() (*RepoMap, error) {
	root := i.Root
	if root == "" {
		root = "."
	}
	m := &RepoMap{
		KeyFiles:  []string{},
		Modules:   []string{},
		Languages: []string{},
		Build:     "",
	}
	langSet := make(map[string]struct{})
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		if strings.Contains(rel, ".git") || strings.Contains(rel, "node_modules") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(p)
		switch ext {
		case ".go":
			langSet["go"] = struct{}{}
			if strings.HasSuffix(rel, "go.mod") || strings.HasSuffix(rel, "main.go") {
				m.KeyFiles = append(m.KeyFiles, rel)
			}
		case ".js", ".ts", ".jsx", ".tsx":
			langSet["javascript"] = struct{}{}
			if strings.Contains(rel, "package.json") {
				m.KeyFiles = append(m.KeyFiles, rel)
			}
		case ".py":
			langSet["python"] = struct{}{}
		case ".rs":
			langSet["rust"] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for lang := range langSet {
		m.Languages = append(m.Languages, lang)
	}
	if _, ok := langSet["go"]; ok {
		m.Build = "go build ./..."
	} else if _, ok := langSet["javascript"]; ok {
		m.Build = "npm install && npm run build"
	}
	return m, nil
}
