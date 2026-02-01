package providers

import "context"

// Message is a chat message (role + content).
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StreamEvent is a single event from a streaming response (token chunk or tool call).
type StreamEvent struct {
	Token    string    `json:"token,omitempty"`
	ToolCall *ToolCall `json:"tool_call,omitempty"`
	Done     bool      `json:"done,omitempty"`
}

// ToolCall represents a tool invocation request from the model.
type ToolCall struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Args string `json:"args,omitempty"`
}

// Provider is the common interface for LLM providers.
type Provider interface {
	// ChatStream streams tokens and optional tool-call events. toolsSchema can be nil for no tools.
	ChatStream(ctx context.Context, messages []Message, toolsSchema interface{}) (Stream, error)
}

// Stream yields token chunks and tool-call events. Caller must call Close when done.
type Stream interface {
	// Next returns the next event. io.EOF when stream is done.
	Next() (StreamEvent, error)
	Close() error
}

// StreamReader adapts a Stream to io.Reader for token-by-token reading (optional).
type StreamReader struct {
	Stream Stream
	buf    []byte
}

func (r *StreamReader) Read(p []byte) (n int, err error) {
	if len(r.buf) > 0 {
		n = copy(p, r.buf)
		r.buf = r.buf[n:]
		return n, nil
	}
	ev, err := r.Stream.Next()
	if err != nil {
		return 0, err
	}
	if ev.Token != "" {
		r.buf = []byte(ev.Token)
		return r.Read(p)
	}
	return r.Read(p)
}
