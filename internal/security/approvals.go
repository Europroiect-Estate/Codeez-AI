package security

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ApprovalScope is the user's choice for a tool call.
type ApprovalScope string

const (
	ScopeDeny    ApprovalScope = "deny"
	ScopeOnce    ApprovalScope = "once"
	ScopeSession ApprovalScope = "session"
	ScopeRepo    ApprovalScope = "repo"
)

// ApprovalPrompt shows agent, tool, reason, args and reads choice from r.
func ApprovalPrompt(r io.Reader, agent, tool, reason string, args interface{}) (ApprovalScope, error) {
	fmt.Printf("[Approval] agent=%s tool=%s reason=%s args=%v\n", agent, tool, reason, args)
	fmt.Print("Deny / Allow once / Allow for session / Always allow for this repo [d/o/s/a]: ")
	sc := bufio.NewScanner(r)
	var line string
	if sc.Scan() {
		line = sc.Text()
	}
	line = strings.TrimSpace(strings.ToLower(line))
	switch line {
	case "d", "deny":
		return ScopeDeny, nil
	case "o", "once":
		return ScopeOnce, nil
	case "s", "session":
		return ScopeSession, nil
	case "a", "repo", "always":
		return ScopeRepo, nil
	default:
		return ScopeDeny, nil
	}
}
