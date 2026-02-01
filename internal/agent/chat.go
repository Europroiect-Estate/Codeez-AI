package agent

import (
	"context"
	"fmt"
	"io"

	"github.com/Europroiect-Estate/Codeez-AI/internal/providers"
	"github.com/Europroiect-Estate/Codeez-AI/internal/store"
	"github.com/Europroiect-Estate/Codeez-AI/internal/ui/tokens"
)

// ChatSession holds session ID and store for persisting messages.
type ChatSession struct {
	Store     *store.Store
	SessionID int64
	Palette   string
}

// ChatOnce sends user message, streams response to w (with palette accent), and persists messages.
func (c *ChatSession) ChatOnce(ctx context.Context, prov providers.Provider, userMsg string, w io.Writer) error {
	_, err := c.Store.AppendMessage(c.SessionID, "user", userMsg)
	if err != nil {
		return err
	}
	msgs, err := c.Store.GetMessages(c.SessionID)
	if err != nil {
		return err
	}
	messages := make([]providers.Message, 0, len(msgs))
	for _, m := range msgs {
		messages = append(messages, providers.Message{Role: m.Role, Content: m.Content})
	}
	stream, err := prov.ChatStream(ctx, messages, nil)
	if err != nil {
		return err
	}
	defer stream.Close()
	palette := tokens.Select(c.Palette)
	acc := palette.AccentANSI()
	reset := "\033[0m"
	var fullContent string
	for {
		ev, err := stream.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if ev.Token != "" {
			fullContent += ev.Token
			fmt.Fprint(w, acc+ev.Token+reset)
		}
	}
	fmt.Fprintln(w)
	_, err = c.Store.AppendMessage(c.SessionID, "assistant", fullContent)
	return err
}
