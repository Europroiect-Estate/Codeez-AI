package agent

// TestAgent decides and runs lint/test/build with approval.
type TestAgent struct {
	Cwd string
}

// SuggestCommand returns a suggested test/lint/build command for the repo.
func (t *TestAgent) SuggestCommand() (string, error) {
	// MVP: common patterns
	return "go test ./...", nil
}
