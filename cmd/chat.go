package main

import (
	"github.com/digkill/telegram-chatgpt/internal/services/chatgpt"
	"github.com/sashabaranov/go-openai"
	"os"
)

func newChat() *chatgpt.ChatGPT {
	var openAIToken = os.Getenv("CHATGPT_TOKEN")
	var openAIURL = os.Getenv("CHATGPT_URL")

	config := openai.DefaultConfig(openAIToken)
	if openAIURL != "" {
		config.BaseURL = openAIURL
	}

	openaiClient := openai.NewClientWithConfig(config)
	chat := chatgpt.NewChatGPT(openaiClient)
	return chat
}
