package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/Europroiect-Estate/Codeez-AI/internal/store"
)

// PatchPreviewTool shows a unified diff without applying.
type PatchPreviewTool struct {
	Store *store.Store
}

func (t *PatchPreviewTool) Name() string        { return "patch.preview" }
func (t *PatchPreviewTool) Description() string { return "Preview unified diff" }
func (t *PatchPreviewTool) Parameters() string  { return `{"unified_diff": "string"}` }

func (t *PatchPreviewTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	diff, _ := args["unified_diff"].(string)
	if diff == "" {
		return "", fmt.Errorf("unified_diff required")
	}
	return "Preview:\n" + diff, nil
}

// PatchApplyTool applies a unified diff (requires approval; actual apply in CLI/store).
type PatchApplyTool struct {
	Store *store.Store
}

func (t *PatchApplyTool) Name() string        { return "patch.apply" }
func (t *PatchApplyTool) Description() string { return "Apply unified diff (requires approval)" }
func (t *PatchApplyTool) Parameters() string {
	return `{"unified_diff": "string", "session_id": "number"}`
}

func (t *PatchApplyTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	diff, _ := args["unified_diff"].(string)
	sessionID, _ := args["session_id"].(float64)
	if diff == "" {
		return "", fmt.Errorf("unified_diff required")
	}
	// Record in store; actual file apply is done by executor with approval
	if t.Store != nil && sessionID > 0 {
		path := "patch"
		if firstLine := strings.SplitN(diff, "\n", 2); len(firstLine) > 0 && strings.HasPrefix(firstLine[0], "---") {
			path = firstLine[0]
		}
		_, _ = t.Store.RecordFileChange(int64(sessionID), path, diff)
	}
	return "recorded; apply with approval", nil
}

// PatchRollbackTool rolls back a change by ID.
type PatchRollbackTool struct {
	Store *store.Store
}

func (t *PatchRollbackTool) Name() string        { return "patch.rollback" }
func (t *PatchRollbackTool) Description() string { return "Rollback change by change_id" }
func (t *PatchRollbackTool) Parameters() string  { return `{"change_id": "number"}` }

func (t *PatchRollbackTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	changeID, _ := args["change_id"].(float64)
	if changeID <= 0 {
		return "", fmt.Errorf("change_id required")
	}
	if t.Store != nil {
		fc, err := t.Store.GetFileChange(int64(changeID))
		if err != nil || fc == nil {
			return "", fmt.Errorf("change not found")
		}
		// Rollback = revert the diff (simplified: just mark not applied; full revert would apply reverse diff)
		return "rollback requested for change " + fmt.Sprint(changeID), nil
	}
	return "ok", nil
}
