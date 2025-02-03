package domains

import "context"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Answer struct {
	Role    string
	Content string
}

type Chat interface {
	Chat(ctx context.Context, messages []Message) (*Answer, error)
}
