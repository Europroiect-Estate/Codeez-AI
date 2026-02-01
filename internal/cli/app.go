package cli

import (
	"os"
	"path/filepath"

	"github.com/Europroiect-Estate/Codeez-AI/internal/app"
	"github.com/spf13/cobra"
)

// getApp loads app with config; repo root is optional (empty if not in repo).
func getApp(cmd *cobra.Command) (*app.App, error) {
	cwd, _ := os.Getwd()
	repoRoot := detectRepoRoot(cwd)
	globalBase := ""
	if home, err := os.UserHomeDir(); err == nil {
		globalBase = home
	}
	return app.New(globalBase, cwd, repoRoot)
}

// detectRepoRoot returns git root or empty string.
func detectRepoRoot(cwd string) string {
	dir := cwd
	for {
		if dir == "" || dir == "/" {
			return ""
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
}
