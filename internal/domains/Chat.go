package domains

import "context"

type Content struct {
	Text  string
	Image string
}

type Message struct {
	Role    string
	Content []Content
}

type Answer struct {
	Role    string
	Content string
}

type Chat interface {
	Chat(ctx context.Context, messages []Message) (*Answer, error)
}
