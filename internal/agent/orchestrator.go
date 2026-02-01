package agent

import (
	"context"
	"fmt"
	"io"

	"github.com/Europroiect-Estate/Codeez-AI/internal/providers"
	"github.com/Europroiect-Estate/Codeez-AI/internal/store"
	"github.com/Europroiect-Estate/Codeez-AI/internal/ui/tokens"
)

// Orchestrator runs the agent loop: plan -> steps -> execute with approvals.
type Orchestrator struct {
	Store     *store.Store
	SessionID int64
	Palette   string
}

// Run runs the task: send to provider with plan prompt, stream response, persist, then print summary.
func (o *Orchestrator) Run(ctx context.Context, prov providers.Provider, task string, w io.Writer) error {
	// Build prompt for planning + execution
	prompt := "Task: " + task + "\n\nFirst output a short numbered plan. Then suggest concrete file changes (as unified diffs if applicable). End with a line: SUGGESTED_COMMIT: <message>"
	sess := &ChatSession{Store: o.Store, SessionID: o.SessionID, Palette: o.Palette}
	if err := sess.ChatOnce(ctx, prov, prompt, w); err != nil {
		return err
	}
	msgs, err := o.Store.GetMessages(o.SessionID)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}
	last := msgs[len(msgs)-1]
	if last.Role != "assistant" {
		return nil
	}
	palette := tokens.Select(o.Palette)
	acc := palette.AccentANSI()
	reset := "\033[0m"
	fmt.Fprintf(w, "\n%s--- Summary ---%s\n", acc, reset)
	fmt.Fprintf(w, "Session %d: %d messages. Review the plan and diffs above.\n", o.SessionID, len(msgs))
	fmt.Fprintf(w, "To apply changes, use: codeez apply (or edit files manually).\n")
	return nil
}
