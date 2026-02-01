package store

import (
	"testing"
)

func TestOpen_and_Migrate_idempotent(t *testing.T) {
	dir := t.TempDir()
	s1, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s1.Close()

	s2, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()
	_ = s2
}

func TestSessions_and_Messages(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	id, err := s.CreateSession("/cwd", "/repo", "ollama", "llama2", "test")
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 {
		t.Errorf("expected positive session id, got %d", id)
	}

	sess, err := s.GetSession(id)
	if err != nil {
		t.Fatal(err)
	}
	if sess == nil || sess.Cwd != "/cwd" || sess.RepoRoot != "/repo" {
		t.Errorf("GetSession: got %+v", sess)
	}

	mid, err := s.AppendMessage(id, "user", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if mid <= 0 {
		t.Errorf("expected positive message id, got %d", mid)
	}

	msgs, err := s.GetMessages(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0].Role != "user" || msgs[0].Content != "hello" {
		t.Errorf("GetMessages: got %+v", msgs)
	}
}

func TestAudit(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	id, err := s.CreateSession("/cwd", "", "ollama", "", "")
	if err != nil {
		t.Fatal(err)
	}

	auditID, err := s.LogToolAudit(id, "planner", "fs.read", map[string]string{"path": "/foo"}, true, "once")
	if err != nil {
		t.Fatal(err)
	}
	if auditID <= 0 {
		t.Errorf("expected positive audit id, got %d", auditID)
	}

	entries, err := s.GetToolAudit(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].ToolName != "fs.read" || !entries[0].Approved {
		t.Errorf("GetToolAudit: got %+v", entries)
	}
}
