package cli

import (
	"fmt"
	"path/filepath"

	"github.com/Europroiect-Estate/Codeez-AI/internal/config"
	"github.com/Europroiect-Estate/Codeez-AI/internal/ui/tokens"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewConfigCmd returns the config command.
func NewConfigCmd() *cobra.Command {
	var global bool
	c := &cobra.Command{
		Use:   "config",
		Short: "Show or edit configuration",
	}
	c.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Print merged config (project overrides global)",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			palette := tokens.Select(app.Config.Palette)
			acc := palette.AccentANSI()
			mut := "\033[38;2;128;128;128m"
			reset := "\033[0m"
			fmt.Printf("%spalette%s = %s%s%s\n", acc, reset, mut, app.Config.Palette, reset)
			fmt.Printf("%sprovider%s = %s%s%s\n", acc, reset, mut, app.Config.Provider, reset)
			fmt.Printf("%smodel%s = %s%s%s\n", acc, reset, mut, app.Config.Model, reset)
			fmt.Printf("%sapi_key%s = %s<redacted>%s\n", acc, reset, mut, reset)
			return nil
		},
	})
	setCmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a config key (e.g. palette corporate)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			key, value := args[0], args[1]
			var configPath string
			if global {
				configPath = filepath.Join(app.GlobalConfigDir(), "config.toml")
			} else {
				if app.ProjectConfigDir() == "" {
					return fmt.Errorf("not in a codeez project; run from a repo with .codeez/ or use --global")
				}
				configPath = filepath.Join(app.ProjectConfigDir(), "config.toml")
			}
			if err := config.EnsureConfigDir(filepath.Dir(configPath)); err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigFile(configPath)
			v.SetConfigType("toml")
			_ = v.ReadInConfig()
			v.Set(key, value)
			if err := v.WriteConfig(); err != nil {
				return err
			}
			fmt.Printf("Set %s = %s\n", key, value)
			return nil
		},
	}
	setCmd.Flags().BoolVar(&global, "global", false, "Write to global config")
	c.AddCommand(setCmd)
	return c
}
