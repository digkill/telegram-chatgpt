package chatGPT

import (
	"context"
	"errors"
	"github.com/digkill/latex2unicode"
	"github.com/sashabaranov/go-openai"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/domains"
)

type ChatGPTComponent struct {
	config *config.ChatGPTConfig
	client *openai.Client
}

func (component *ChatGPTComponent) Init() bool {
	var openAIToken = component.config.Token

	// URL HTTP-прокси
	//	proxyURL, err := url.Parse("https://168.228.44.66:999")
	//	if err != nil {
	//		fmt.Println("😡😡😡😡😡")
	//		fmt.Println("Ошибка парсинга прокси:", err)
	//		fmt.Println("😡😡😡😡😡")
	//	}

	//	transport := &http.Transport{
	//		Proxy: http.ProxyURL(proxyURL),
	//	}
	//	client := &http.Client{Transport: transport}

	clientConfig := openai.DefaultConfig(openAIToken)
	//	clientConfig.HTTPClient = client
	component.client = openai.NewClientWithConfig(clientConfig)

	return true
}

type Option func(*ChatGPTComponent)

func SetModel(model string) Option {
	return func(gpt *ChatGPTComponent) {
		gpt.config.Model = model
	}
}

func (c ChatGPTComponent) GetChat() *openai.Client {
	return c.client
}

func (c ChatGPTComponent) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (*domains.Answer, error) {

	//chatGPTMessages := c.makeChatGPTMessage(messages)
	response, err := c.send(ctx, messages)
	if err != nil {
		return nil, err
	}
	responseLatex := domains.Answer{
		Role:    response.Role,
		Content: latex2unicode.ConvertLatexToUnicode(response.Content),
	}

	return &responseLatex, nil
	// return response, nil
}

func (c ChatGPTComponent) makeChatGPTMessage(messages []domains.Message) []openai.ChatCompletionMessage {

	chatGPTMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, message := range messages {
		chatGPTMessages = append(chatGPTMessages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	return chatGPTMessages
}

func (c ChatGPTComponent) send(ctx context.Context, chatGPTMessages []openai.ChatCompletionMessage) (*domains.Answer, error) {

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    c.config.Model,
			Messages: chatGPTMessages,
			// MaxTokens: 100,
			// Store: false,
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

func (c ChatGPTComponent) convertAnswer(openaiResp openai.ChatCompletionResponse) *domains.Answer {

	choices := openaiResp.Choices[0]

	return &domains.Answer{
		Role:    choices.Message.Role,
		Content: choices.Message.Content,
	}
}

func NewChatGPT(config *config.ChatGPTConfig) *ChatGPTComponent {
	chatGPTComponent := &ChatGPTComponent{
		config: config,
	}

	chatGPTComponent.Init()
	return chatGPTComponent
}
