package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/Europroiect-Estate/Codeez-AI/internal/agent"
	"github.com/Europroiect-Estate/Codeez-AI/internal/providers"
	"github.com/Europroiect-Estate/Codeez-AI/internal/store"
	"github.com/Europroiect-Estate/Codeez-AI/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// NewChatCmd returns the chat command.
func NewChatCmd() *cobra.Command {
	var noTUI bool
	cmd := &cobra.Command{
		Use:   "chat [message]",
		Short: "Interactive TUI or streaming chat session",
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
			sessionID, err := st.CreateSession(app.Cwd, app.RepoRoot, app.Config.Provider, app.Config.Model, "chat")
			if err != nil {
				return err
			}
			prov, err := providers.NewFactory(app.Config).Get()
			if err != nil {
				return err
			}
			chatSess := &agent.ChatSession{Store: st, SessionID: sessionID, Palette: app.Config.Palette}
			ctx := context.Background()
			var userMsg string
			if len(args) > 0 {
				userMsg = args[0]
			} else if !noTUI && isTTY(os.Stdin) && isTTY(os.Stdout) {
				// Run TUI
				p := tea.NewProgram(ui.NewModel(app.Config.Palette), tea.WithAltScreen())
				if _, err := p.Run(); err != nil {
					return err
				}
				return nil
			} else {
				if noTUI || !isTTY(os.Stdin) {
					sc := bufio.NewScanner(os.Stdin)
					if sc.Scan() {
						userMsg = sc.Text()
					}
				} else {
					fmt.Print("You: ")
					sc := bufio.NewScanner(os.Stdin)
					if sc.Scan() {
						userMsg = sc.Text()
					}
				}
			}
			if userMsg == "" {
				return nil
			}
			return chatSess.ChatOnce(ctx, prov, userMsg, os.Stdout)
		},
	}
	cmd.Flags().BoolVar(&noTUI, "no-tui", false, "Stream to terminal (no TUI)")
	return cmd
}

func isTTY(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
