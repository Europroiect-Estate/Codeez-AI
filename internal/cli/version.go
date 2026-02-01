package cli

import (
	"fmt"

	"github.com/Europroiect-Estate/Codeez-AI/internal/ui/tokens"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

// NewVersionCmd returns the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print codeez version",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			palette := tokens.Select(app.Config.Palette)
			acc := palette.AccentANSI()
			reset := "\033[0m"
			fmt.Printf("codeez version %s%s%s\n", acc, Version, reset)
			return nil
		},
	}
}
