package agent

// GitAgent handles git workflows: branch, stage, commit message, log summary; avoids committing secrets.
type GitAgent struct{}

// SuggestCommitMessage returns a suggested commit message for the given diff/summary.
func (g *GitAgent) SuggestCommitMessage(summary string) (string, error) {
	return "feat: " + summary, nil
}
