package agent

// CoderAgent produces patches/diffs. Used by orchestrator when executing a step.
type CoderAgent struct{}

// SuggestPatch returns a unified diff for the given step. For MVP this is done via LLM in orchestrator.
func (c *CoderAgent) SuggestPatch(step, context string) (string, error) {
	return "", nil
}
