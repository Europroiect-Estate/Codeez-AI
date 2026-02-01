package security

import "testing"

func TestRedact(t *testing.T) {
	in := "api_key=sk-12345678901234567890"
	out := Redact(in)
	if out == in {
		t.Error("expected redaction")
	}
	if out != "api_key=***" && out != "api_key=sk-***" {
		t.Errorf("got %q", out)
	}
}

func TestContainsSecret(t *testing.T) {
	if !ContainsSecret("Bearer abc123token") {
		t.Error("expected true for Bearer token")
	}
	if ContainsSecret("hello world") {
		t.Error("expected false for plain text")
	}
}
