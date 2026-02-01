package cli

import (
	"fmt"

	"github.com/Europroiect-Estate/Codeez-AI/internal/index"
	"github.com/spf13/cobra"
)

// NewIndexCmd returns the index command.
func NewIndexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "index",
		Short: "Build repo map (key files, languages, build hints)",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			dir := app.ProjectConfigDir()
			if dir == "" {
				dir = app.Cwd + "/.codeez"
			}
			ix := index.Indexer{Root: app.Cwd}
			m, err := ix.Index()
			if err != nil {
				return err
			}
			if err := m.Save(dir); err != nil {
				return err
			}
			fmt.Printf("Indexed: %d key files, languages %v\n", len(m.KeyFiles), m.Languages)
			return nil
		},
	}
}
