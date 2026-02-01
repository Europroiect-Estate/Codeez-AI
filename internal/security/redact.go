package security

import "regexp"

// Redact replaces known secret patterns with a placeholder.
func Redact(s string) string {
	out := s
	// Common API key / token patterns (simplified)
	patterns := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`(?i)(sk-[a-zA-Z0-9]{20,})`), `sk-***`},
		{regexp.MustCompile(`(?i)(api[_-]?key\s*[:=]\s*)["']?[a-zA-Z0-9_\-]{20,}`), `$1***`},
		{regexp.MustCompile(`(?i)(Bearer\s+)[a-zA-Z0-9_\-\.]+`), `$1***`},
		{regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`), `ghp_***`},
		{regexp.MustCompile(`AKIA[0-9A-Z]{16}`), `AKIA***`},
	}
	for _, p := range patterns {
		out = p.re.ReplaceAllString(out, p.repl)
	}
	return out
}

// RedactEntropy replaces long high-entropy strings (e.g. tokens) with ***.
func RedactEntropy(s string, minLen int) string {
	if minLen <= 0 {
		minLen = 32
	}
	// Replace long alphanumeric stretches (likely tokens)
	tokenRe := regexp.MustCompile(`[A-Za-z0-9+/=_-]{32,}`)
	return tokenRe.ReplaceAllString(s, "***")
}

// RedactForLog returns a string safe for logging (redact known secrets).
func RedactForLog(s string) string {
	return Redact(s)
}

// ContainsSecret returns true if the string likely contains a secret (for blocking commit).
func ContainsSecret(s string) bool {
	redacted := Redact(s)
	return redacted != s
}
