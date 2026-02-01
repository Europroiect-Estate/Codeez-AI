package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Europroiect-Estate/Codeez-AI/internal/agent"
	"github.com/Europroiect-Estate/Codeez-AI/internal/providers"
	"github.com/Europroiect-Estate/Codeez-AI/internal/store"
	"github.com/spf13/cobra"
)

// NewRunCmd returns the run command.
func NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run [task]",
		Short: "Non-TUI agentic task execution with progress and approvals",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			dir := app.ProjectConfigDir()
			if dir == "" {
				dir = app.GlobalConfigDir()
			}
			st, err := store.Open(dir)
			if err != nil {
				return err
			}
			defer st.Close()
			task := strings.Join(args, " ")
			sessionID, err := st.CreateSession(app.Cwd, app.RepoRoot, app.Config.Provider, app.Config.Model, task)
			if err != nil {
				return err
			}
			prov, err := providers.NewFactory(app.Config).Get()
			if err != nil {
				return err
			}
			orch := &agent.Orchestrator{Store: st, SessionID: sessionID, Palette: app.Config.Palette}
			ctx := context.Background()
			if err := orch.Run(ctx, prov, task, os.Stdout); err != nil {
				return fmt.Errorf("run: %w", err)
			}
			return nil
		},
	}
}
