package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Europroiect-Estate/Codeez-AI/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewProviderCmd returns the provider command.
func NewProviderCmd() *cobra.Command {
	providerCmd := &cobra.Command{
		Use:   "provider",
		Short: "List, set, and test LLM providers",
	}
	providerCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available providers",
		RunE:  runProviderList,
	})
	providerCmd.AddCommand(&cobra.Command{
		Use:   "set [name]",
		Short: "Set active provider (ollama, openai, anthropic)",
		Args:  cobra.ExactArgs(1),
		RunE:  runProviderSet,
	})
	providerCmd.AddCommand(&cobra.Command{
		Use:   "set-key [name] [key]",
		Short: "Set API key for a provider",
		Args:  cobra.ExactArgs(2),
		RunE:  runProviderSetKey,
	})
	providerCmd.AddCommand(&cobra.Command{
		Use:   "set-model [name] [model]",
		Short: "Set model for a provider",
		Args:  cobra.ExactArgs(2),
		RunE:  runProviderSetModel,
	})
	providerCmd.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Test connectivity to active provider",
		RunE:  runProviderTest,
	})
	return providerCmd
}

func runProviderList(cmd *cobra.Command, args []string) error {
	fmt.Println("ollama")
	fmt.Println("openai")
	fmt.Println("anthropic")
	return nil
}

func runProviderSet(cmd *cobra.Command, args []string) error {
	name := strings.ToLower(args[0])
	for _, p := range config.ValidProviders() {
		if p == name {
			return writeConfigKey(cmd, "provider", name)
		}
	}
	return fmt.Errorf("unknown provider %q", name)
}

func runProviderSetKey(cmd *cobra.Command, args []string) error {
	name, key := strings.ToLower(args[0]), args[1]
	for _, p := range config.ValidProviders() {
		if p == name {
			return writeConfigKey(cmd, "providers."+name+".api_key", key)
		}
	}
	return fmt.Errorf("unknown provider %q", name)
}

func runProviderSetModel(cmd *cobra.Command, args []string) error {
	name, model := strings.ToLower(args[0]), args[1]
	for _, p := range config.ValidProviders() {
		if p == name {
			return writeConfigKey(cmd, "providers."+name+".model", model)
		}
	}
	return fmt.Errorf("unknown provider %q", name)
}

func writeConfigKey(cmd *cobra.Command, key, value string) error {
	app, err := getApp(cmd)
	if err != nil {
		return err
	}
	configPath := filepath.Join(app.GlobalConfigDir(), "config.toml")
	if err := config.EnsureConfigDir(app.GlobalConfigDir()); err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("toml")
	_ = v.ReadInConfig()
	v.Set(key, value)
	return v.WriteConfig()
}

func runProviderTest(cmd *cobra.Command, args []string) error {
	app, err := getApp(cmd)
	if err != nil {
		return err
	}
	provider := app.Config.Provider
	switch provider {
	case config.ProviderOllama:
		resp, err := http.Get("http://localhost:11434/api/tags")
		if err != nil {
			return fmt.Errorf("ollama: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ollama returned %d", resp.StatusCode)
		}
		fmt.Println("Ollama: OK")
	case config.ProviderOpenAI, config.ProviderAnthropic:
		fmt.Printf("%s: set api_key and run test (stub)\n", provider)
	default:
		return fmt.Errorf("no provider set; run codeez provider set ollama")
	}
	return nil
}

// readLine reads one line from stdin (for approval flows later).
func readLine() (string, error) {
	sc := bufio.NewScanner(os.Stdin)
	if sc.Scan() {
		return strings.TrimSpace(sc.Text()), nil
	}
	return "", sc.Err()
}
