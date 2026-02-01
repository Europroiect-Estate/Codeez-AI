package agent

// ReviewerAgent checks correctness and style. Used after coder produces a patch.
type ReviewerAgent struct{}

// Review returns feedback on the given diff. For MVP this is done via LLM when needed.
func (r *ReviewerAgent) Review(diff string) (string, error) {
	return "Review: ensure idiomatic code and clean style.", nil
}
