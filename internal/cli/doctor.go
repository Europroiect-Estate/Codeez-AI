package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Europroiect-Estate/Codeez-AI/internal/ui/tokens"
	"github.com/spf13/cobra"
)

// NewDoctorCmd returns the doctor command.
func NewDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose environment: git, rg, ollama, node, permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}
			palette := tokens.Select(app.Config.Palette)
			okStyle := palette.AccentANSI()
			reset := "\033[0m"
			muted := "\033[90m"

			fmt.Printf("%s  codeez doctor — verificare mediu%s\n\n", okStyle, reset)
			checks := []struct {
				name string
				fn   func() string
			}{
				{"git", checkGit},
				{"rg (ripgrep)", checkRg},
				{"ollama", checkOllama},
				{"node", checkNode},
				{"~/.config/codeez writable", func() string { return checkDirWritable(app.GlobalConfigDir()) }},
				{".codeez writable", func() string {
					if app.ProjectConfigDir() == "" {
						return "skip (no project)"
					}
					return checkDirWritable(app.ProjectConfigDir())
				}},
			}

			fmt.Println(muted + "  Check                    Status" + reset)
			fmt.Println(muted + "  ─────────────────────────────────" + reset)
			for _, c := range checks {
				status := c.fn()
				if status == "OK" {
					fmt.Printf("  %-24s %s✓ OK%s\n", c.name, okStyle, reset)
				} else {
					fmt.Printf("  %-24s %s\n", c.name, status)
				}
			}
			fmt.Printf("\n%s  Sfat: codeez provider set ollama && codeez init%s\n", muted, reset)
			return nil
		},
	}
}

func checkGit() string {
	_, err := exec.LookPath("git")
	if err != nil {
		return "Missing"
	}
	return "OK"
}

func checkRg() string {
	_, err := exec.LookPath("rg")
	if err != nil {
		return "Missing"
	}
	return "OK"
}

func checkOllama() string {
	_, err := exec.LookPath("ollama")
	if err != nil {
		return "Missing"
	}
	return "OK"
}

func checkNode() string {
	_, err := exec.LookPath("node")
	if err != nil {
		return "Missing (optional)"
	}
	return "OK"
}

func checkDirWritable(dir string) string {
	if dir == "" {
		return "skip"
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "Error: " + err.Error()
	}
	path := filepath.Join(dir, ".write-test")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		return "Not writable"
	}
	_ = os.Remove(path)
	return "OK"
}
