package chatgpt

import (
	"context"
	"errors"
	"github.com/digkill/telegram-chatgpt/internal/domains"
	"github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	model  string
	client *openai.Client
}

type Option func(*ChatGPT)

func SetModel(model string) Option {
	return func(gpt *ChatGPT) {
		gpt.model = model
	}
}

func NewChatGPT(client *openai.Client, opts ...Option) *ChatGPT {

	chatGPT := &ChatGPT{
		model:  openai.GPT4o,
		client: client,
	}

	for _, opt := range opts {
		opt(chatGPT)
	}

	return chatGPT
}

func (c ChatGPT) Chat(ctx context.Context, messages []domains.Message) (*domains.Answer, error) {

	chatGPTMessages := c.makeChatGPTMessage(messages)

	return c.send(ctx, chatGPTMessages)
}

func (c ChatGPT) makeChatGPTMessage(messages []domains.Message) []openai.ChatCompletionMessage {

	chatGPTMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, message := range messages {
		chatGPTMessages = append(chatGPTMessages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	return chatGPTMessages
}

func (c ChatGPT) send(ctx context.Context, chatGPTMessages []openai.ChatCompletionMessage) (*domains.Answer, error) {

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: chatGPTMessages,
	})
	if err != nil {
		return nil, err
	}

	if choices := resp.Choices; len(choices) == 0 {
		return nil, errors.New("got empty ChatGPT response")
	}

	answer := c.convertAnswer(resp)
	return answer, nil
}

func (c ChatGPT) convertAnswer(openaiResp openai.ChatCompletionResponse) *domains.Answer {

	choices := openaiResp.Choices[0]

	return &domains.Answer{
		Role:    choices.Message.Role,
		Content: choices.Message.Content,
	}
}
