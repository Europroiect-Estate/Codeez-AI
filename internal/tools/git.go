package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Europroiect-Estate/Codeez-AI/internal/security"
)

// GitTool wraps system git. All operations execute `git` with approval.
type GitTool struct {
	Cwd     string
	Sandbox *security.Sandbox
}

func (t *GitTool) Name() string { return "git" }
func (t *GitTool) Description() string {
	return "Git operations (root, status, diff, add, commit, branch, checkout, log, remote, push, pull)"
}
func (t *GitTool) Parameters() string { return `{"subcommand": "string", "args": ["string"]}` }

func (t *GitTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	sub, _ := args["subcommand"].(string)
	if sub == "" {
		return "", fmt.Errorf("subcommand required")
	}
	cwd := t.Cwd
	if cwd == "" {
		cwd = "."
	}
	var argList []string
	if a, ok := args["args"].([]interface{}); ok {
		for _, v := range a {
			if s, ok := v.(string); ok {
				argList = append(argList, s)
			}
		}
	}
	return runGit(ctx, cwd, sub, argList)
}

func runGit(ctx context.Context, cwd, sub string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{sub}, args...)...)
	cmd.Dir = cwd
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("git %s: %w", sub, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Helpers used by CLI git subcommands (same logic, explicit subcommand names).
func GitRoot(cwd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func GitStatus(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "status", nil)
}

func GitDiff(cwd string, paths ...string) (string, error) {
	return runGit(context.Background(), cwd, "diff", paths)
}

func GitAdd(cwd string, paths ...string) (string, error) {
	if len(paths) == 0 {
		return runGit(context.Background(), cwd, "add", []string{"--all"})
	}
	return runGit(context.Background(), cwd, "add", paths)
}

func GitCommit(cwd, message string) (string, error) {
	return runGit(context.Background(), cwd, "commit", []string{"-m", message})
}

// GitStagedDiff returns the staged diff (for secret check).
func GitStagedDiff(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "diff", []string{"--cached"})
}

func GitBranch(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "branch", nil)
}

func GitCheckout(cwd, branch string) (string, error) {
	return runGit(context.Background(), cwd, "checkout", []string{branch})
}

func GitLog(cwd string, oneline bool, max int) (string, error) {
	args := []string{"log"}
	if oneline {
		args = append(args, "--oneline")
	}
	if max > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", max))
	}
	return runGit(context.Background(), cwd, "log", args)
}

func GitRemote(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "remote", []string{"-v"})
}

func GitPush(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "push", nil)
}

func GitPull(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "pull", nil)
}

func GitInit(cwd string) (string, error) {
	return runGit(context.Background(), cwd, "init", nil)
}

func IsRepo(cwd string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = cwd
	err := cmd.Run()
	return err == nil
}
