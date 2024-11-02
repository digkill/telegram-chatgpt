package main

import (
	"github.com/digkill/telegram-chatgpt/internal/domains"
	"github.com/digkill/telegram-chatgpt/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/pkg/escape"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(); err != nil {
		logrus.Warnf("load env failed: %v", err)
	}

	var (
		port = os.Getenv("HTTP_PORT")
	)

	if port == "" {
		port = ":8080"
	}

	var chat domains.Chat
	{
		chat = newChat()
	}

	r := server.NewHTTPServer()
	r.POST("/message", func(ctx *gin.Context) {

		msg, _ := ctx.GetPostForm("message")

		messages := makeMessages(msg)

		answer, err := chat.Chat(ctx, messages)
		if err != nil {
			logrus.Errorf("ChatGPT: chat error: %v", err)
		}

		logrus.Debugf("ChatGPT: chat: message=%s answer: %s", msg, answer.Content)

		result := escape.String(answer.Content)
		ctx.String(http.StatusOK, result)
	})

	logrus.Fatal(r.Run(port))
}

func makeMessages(content string) []domains.Message {

	return []domains.Message{
		{
			Role:    "user",
			Content: content,
		},
	}
}
