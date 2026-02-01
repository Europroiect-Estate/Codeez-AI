package tools

import "context"

// Tool is the interface for agent tools (fs, patch, git, cmd, repo).
type Tool interface {
	Name() string
	Description() string
	Parameters() string // JSON schema or description
	Execute(ctx context.Context, args map[string]interface{}) (result string, err error)
}
