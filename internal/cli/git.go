package cli

import (
	"fmt"
	"os"

	"github.com/Europroiect-Estate/Codeez-AI/internal/agent"
	"github.com/Europroiect-Estate/Codeez-AI/internal/security"
	"github.com/Europroiect-Estate/Codeez-AI/internal/tools"
	"github.com/spf13/cobra"
)

// NewGitCmd returns the git command (subcommands: init, status, diff, commit, branch, checkout, add, log, remote, push, pull).
func NewGitCmd() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:   "git",
		Short: "Explicit git operations",
	}
	gitCmd.AddCommand(newGitSub("init", "Init repo", nil, runGitInit))
	gitCmd.AddCommand(newGitSub("status", "Git status", nil, runGitStatus))
	gitCmd.AddCommand(newGitSub("diff", "Git diff", nil, runGitDiff))
	gitCmd.AddCommand(newGitSub("add", "Git add", cobra.ArbitraryArgs, runGitAdd))
	gitCmd.AddCommand(gitCommitCmd())
	gitCmd.AddCommand(newGitSub("branch", "Git branch", nil, runGitBranch))
	gitCmd.AddCommand(newGitSub("checkout", "Git checkout", cobra.ExactArgs(1), runGitCheckout))
	gitCmd.AddCommand(gitLogCmd())
	gitCmd.AddCommand(newGitSub("remote", "Git remote", nil, runGitRemote))
	gitCmd.AddCommand(newGitSub("push", "Git push", nil, runGitPush))
	gitCmd.AddCommand(newGitSub("pull", "Git pull", nil, runGitPull))
	return gitCmd
}

func newGitSub(use, short string, args cobra.PositionalArgs, runE func(*cobra.Command, []string) error) *cobra.Command {
	c := &cobra.Command{Use: use, Short: short, RunE: runE}
	if args != nil {
		c.Args = args
	}
	return c
}

func gitLogCmd() *cobra.Command {
	c := &cobra.Command{Use: "log", Short: "Git log", RunE: runGitLog}
	c.Flags().Bool("oneline", false, "One line per commit")
	c.Flags().Int("max", 20, "Max commits")
	return c
}

func getCwd(cmd *cobra.Command) string {
	app, err := getApp(cmd)
	if err != nil {
		return "."
	}
	return app.Cwd
}

func runGitInit(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	out, err := tools.GitInit(cwd)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitStatus(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	out, err := tools.GitStatus(cwd)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitDiff(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	out, err := tools.GitDiff(cwd, args...)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitAdd(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	var paths []string
	for _, a := range args {
		if a == "--all" || a == "-A" {
			paths = nil
			break
		}
		paths = append(paths, a)
	}
	out, err := tools.GitAdd(cwd, paths...)
	if err != nil {
		return err
	}
	if out != "" {
		fmt.Println(out)
	}
	return nil
}

func gitCommitCmd() *cobra.Command {
	c := &cobra.Command{Use: "commit", Short: "Git commit", RunE: runGitCommit}
	c.Flags().StringP("message", "m", "", "Commit message")
	c.Args = cobra.NoArgs
	return c
}

func runGitCommit(cmd *cobra.Command, args []string) error {
	msg, _ := cmd.Flags().GetString("message")
	if msg == "" {
		return fmt.Errorf("commit message required (-m \"...\")")
	}
	cwd := getCwd(cmd)
	// Security: check staged content for secrets
	sec := &agent.SecurityAgent{}
	diff, err := tools.GitStagedDiff(cwd)
	if err == nil && diff != "" {
		if sec.BlockCommit(diff) {
			fmt.Print("Detected secrets in staged content. Type 'yes' to override: ")
			var override string
			_, _ = fmt.Scanln(&override)
			if override != "yes" {
				fmt.Println("Commit aborted.")
				return nil
			}
		}
	}
	out, err := tools.GitCommit(cwd, msg)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitBranch(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	out, err := tools.GitBranch(cwd)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitCheckout(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	out, err := tools.GitCheckout(cwd, args[0])
	if err != nil {
		return err
	}
	if out != "" {
		fmt.Println(out)
	}
	return nil
}

func runGitLog(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	oneline, _ := cmd.Flags().GetBool("oneline")
	max, _ := cmd.Flags().GetInt("max")
	if max <= 0 {
		max = 20
	}
	out, err := tools.GitLog(cwd, oneline, max)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitRemote(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	out, err := tools.GitRemote(cwd)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func runGitPush(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	scope, err := security.ApprovalPrompt(os.Stdin, "user", "git.push", "network write", nil)
	if err != nil {
		return err
	}
	if scope == security.ScopeDeny {
		fmt.Println("Push denied.")
		return nil
	}
	out, err := tools.GitPush(cwd)
	if err != nil {
		return err
	}
	if out != "" {
		fmt.Println(out)
	}
	return nil
}

func runGitPull(cmd *cobra.Command, args []string) error {
	cwd := getCwd(cmd)
	scope, err := security.ApprovalPrompt(os.Stdin, "user", "git.pull", "network read", nil)
	if err != nil {
		return err
	}
	if scope == security.ScopeDeny {
		fmt.Println("Pull denied.")
		return nil
	}
	out, err := tools.GitPull(cwd)
	if err != nil {
		return err
	}
	if out != "" {
		fmt.Println(out)
	}
	return nil
}
