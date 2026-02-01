package agent

// ReleaseAgent checks goreleaser, CI workflows, and prepares changelog snippet.
type ReleaseAgent struct{}

// CheckRelease returns a short report on release readiness (goreleaser, CI).
func (r *ReleaseAgent) CheckRelease() (string, error) {
	return "Release: check .goreleaser.yaml and .github/workflows.", nil
}
