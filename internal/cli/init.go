package cli

import (
	"fmt"
	"path/filepath"

	"github.com/Europroiect-Estate/Codeez-AI/internal/config"
	"github.com/Europroiect-Estate/Codeez-AI/internal/store"
	"github.com/spf13/cobra"
)

// NewInitCmd returns the init command.
func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize project (.codeez/ config) and optionally detect git root",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			dir := app.ProjectConfigDir()
			if dir == "" {
				dir = filepath.Join(app.Cwd, ".codeez")
			}
			if err := config.EnsureConfigDir(dir); err != nil {
				return err
			}
			s, err := store.Open(dir)
			if err != nil {
				return err
			}
			_ = s.Close()
			fmt.Println("Initialized .codeez at", dir)
			return nil
		},
	}
}
