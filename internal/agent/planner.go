package agent

// PlannerAgent turns a task into steps and tool suggestions.
// Used by the orchestrator to get a plan before executing.
type PlannerAgent struct{}

// Plan returns a short plan (steps) for the task. For MVP this is done via LLM in orchestrator.
func (p *PlannerAgent) Plan(task string) ([]string, error) {
	// Placeholder: real implementation uses provider with "output numbered steps"
	return []string{"1. Analyze task", "2. Propose changes", "3. Apply with approval"}, nil
}
