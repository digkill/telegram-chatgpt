package domains

import "context"

type Message struct {
	Role    string
	Content string
}

type Answer struct {
	Role    string
	Content string
}

type Chat interface {
	Chat(ctx context.Context, messages []Message) (*Answer, error)
}
