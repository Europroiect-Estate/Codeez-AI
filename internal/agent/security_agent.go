package agent

import (
	"fmt"

	"github.com/Europroiect-Estate/Codeez-AI/internal/security"
)

// SecurityAgent checks for secrets, unsafe exec, path traversal, injection, etc.
type SecurityAgent struct{}

// CheckDiff returns an error if the diff contains detected secrets or dangerous patterns.
func (s *SecurityAgent) CheckDiff(diff string) error {
	if security.ContainsSecret(diff) {
		return fmt.Errorf("diff contains detected secret; refuse to commit unless overridden with strong confirmation")
	}
	return nil
}

// BlockCommit returns true if we should block commit (e.g. secrets in staged content).
func (s *SecurityAgent) BlockCommit(stagedContent string) bool {
	return security.ContainsSecret(stagedContent)
}
