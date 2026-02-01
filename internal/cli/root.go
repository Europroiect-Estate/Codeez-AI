package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root cobra command.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "codeez",
		Short: "Codeez by Europroiect Estate — agentic coding CLI",
		Long: `Codeez — CLI agentic pentru codare: rapid, sigur, multi-provider (Ollama, OpenAI, Anthropic).
TUI opțional, aprobări explicite, fără telemetrie. Rulează local pe repository-ul tău.

Instalare:
  curl -sSL https://raw.githubusercontent.com/Europroiect-Estate/Codeez-AI/main/scripts/install.sh | bash

Sfat: rulează 'codeez completion' pentru completare în shell (bash/zsh/fish).`,
		Example: `  codeez init
  codeez provider set ollama
  codeez doctor
  codeez chat --no-tui
  codeez run "adaugă o comandă hello"`,
	}
	root.AddCommand(NewInitCmd())
	root.AddCommand(NewChatCmd())
	root.AddCommand(NewRunCmd())
	root.AddCommand(NewProviderCmd())
	root.AddCommand(NewDoctorCmd())
	root.AddCommand(NewIndexCmd())
	root.AddCommand(NewApplyCmd())
	root.AddCommand(NewConfigCmd())
	root.AddCommand(NewGitCmd())
	root.AddCommand(NewVersionCmd())
	return root
}
