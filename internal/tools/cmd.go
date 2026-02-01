package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Denylist of dangerous command prefixes (always require explicit approval; can still be denied).
var cmdDenylist = []string{
	"rm -rf", "rm -fr", "mkfs", "dd if=", "curl | sh", "wget -O- | sh",
	":(){ :|:& };:", "chmod 777", "chown ", "> /dev/sd",
}

const (
	defaultCmdTimeout = 30 * time.Second
	maxOutputSize     = 512 * 1024 // 512KB
)

// CmdRunTool runs a command (always requires approval).
type CmdRunTool struct {
	Timeout time.Duration
}

func (t *CmdRunTool) Name() string        { return "cmd.run" }
func (t *CmdRunTool) Description() string { return "Run command (always requires approval)" }
func (t *CmdRunTool) Parameters() string {
	return `{"command": "string", "args": ["string"], "cwd": "string"}`
}

func (t *CmdRunTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	command, _ := args["command"].(string)
	if command == "" {
		return "", fmt.Errorf("command required")
	}
	cmdLine := command
	if a, ok := args["args"].([]interface{}); ok {
		for _, v := range a {
			if s, ok := v.(string); ok {
				cmdLine += " " + s
			}
		}
	}
	for _, d := range cmdDenylist {
		if strings.Contains(cmdLine, d) {
			return "", fmt.Errorf("command denied by policy: %s", d)
		}
	}
	timeout := t.Timeout
	if timeout == 0 {
		timeout = defaultCmdTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cwd, _ := args["cwd"].(string)
	cmd := exec.CommandContext(ctx, command)
	if args["args"] != nil {
		if a, ok := args["args"].([]interface{}); ok {
			for _, v := range a {
				if s, ok := v.(string); ok {
					cmd.Args = append(cmd.Args, s)
				}
			}
		}
	}
	if cwd != "" {
		cmd.Dir = cwd
	}
	out, err := cmd.CombinedOutput()
	if len(out) > maxOutputSize {
		out = append(out[:maxOutputSize], []byte("\n... (truncated)")...)
	}
	if err != nil {
		return string(out), fmt.Errorf("%w: %s", err, out)
	}
	return string(out), nil
}
