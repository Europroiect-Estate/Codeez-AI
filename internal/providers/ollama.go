package providers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultOllamaBase = "http://localhost:11434"

// OllamaConfig configures the Ollama provider.
type OllamaConfig struct {
	BaseURL string
	Model   string
}

// Ollama implements Provider for Ollama (local).
type Ollama struct {
	client *http.Client
	base   string
	model  string
}

// NewOllama returns an Ollama provider.
func NewOllama(cfg OllamaConfig) *Ollama {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultOllamaBase
	}
	return &Ollama{
		client: &http.Client{Timeout: 120 * time.Second},
		base:   strings.TrimSuffix(cfg.BaseURL, "/"),
		model:  cfg.Model,
	}
}

// ollamaChatReq matches Ollama POST /api/chat request.
type ollamaChatReq struct {
	Model    string      `json:"model"`
	Messages []ollamaMsg `json:"messages"`
	Stream   bool        `json:"stream"`
}

type ollamaMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ollamaChunk is one NDJSON line from Ollama stream.
type ollamaChunk struct {
	Message struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	} `json:"message"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
}

// ChatStream streams from Ollama /api/chat.
func (o *Ollama) ChatStream(ctx context.Context, messages []Message, toolsSchema interface{}) (Stream, error) {
	model := o.model
	if model == "" {
		model = "llama3.2"
	}
	msgs := make([]ollamaMsg, len(messages))
	for i, m := range messages {
		msgs[i] = ollamaMsg(m)
	}
	body := ollamaChatReq{Model: model, Messages: msgs, Stream: true}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.base+"/api/chat", strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("ollama: %s", resp.Status)
	}
	return &ollamaStream{body: resp.Body}, nil
}

type ollamaStream struct {
	body    io.ReadCloser
	scan    *bufio.Scanner
	lastLen int // for delta-only emission (Ollama sends full content per chunk)
}

func (s *ollamaStream) Next() (StreamEvent, error) {
	if s.body == nil {
		return StreamEvent{}, io.EOF
	}
	if s.scan == nil {
		s.scan = bufio.NewScanner(s.body)
	}
	if !s.scan.Scan() {
		_ = s.body.Close()
		s.body = nil
		if err := s.scan.Err(); err != nil {
			return StreamEvent{}, err
		}
		return StreamEvent{}, io.EOF
	}
	line := s.scan.Text()
	var c ollamaChunk
	if err := json.Unmarshal([]byte(line), &c); err != nil {
		return StreamEvent{}, err
	}
	ev := StreamEvent{}
	content := c.Message.Content
	if len(content) > s.lastLen {
		ev.Token = content[s.lastLen:]
		s.lastLen = len(content)
	}
	if c.Done {
		ev.Done = true
	}
	return ev, nil
}

func (s *ollamaStream) Close() error {
	if s.body != nil {
		return s.body.Close()
	}
	return nil
}
