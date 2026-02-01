package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Europroiect-Estate/Codeez-AI/internal/config"
)

func TestOllama_ChatStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" || r.Method != http.MethodPost {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/x-ndjson")
		// Simulate NDJSON stream (Ollama sends cumulative content per chunk)
		w.Write([]byte(`{"message":{"content":"Hello","role":"assistant"},"done":false}
{"message":{"content":"Hello world","role":"assistant"},"done":true}
`))
	}))
	defer server.Close()

	o := NewOllama(OllamaConfig{BaseURL: server.URL, Model: "test"})
	ctx := context.Background()
	stream, err := o.ChatStream(ctx, []Message{{Role: "user", Content: "hi"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var full string
	for {
		ev, err := stream.Next()
		if err != nil {
			break
		}
		full += ev.Token
	}
	if full != "Hello world" {
		t.Errorf("got %q", full)
	}
}

func TestFactory_Get_ollama(t *testing.T) {
	f := NewFactory(&config.Config{Provider: "ollama", Model: "llama3.2"})
	prov, err := f.Get()
	if err != nil {
		t.Fatal(err)
	}
	if prov == nil {
		t.Fatal("expected non-nil provider")
	}
	_ = prov
}
