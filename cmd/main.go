package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(); err != nil {
		logrus.Warnf("load env failed: %v", err)
	}
	runTelegram()
	/*var (
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

			actionInfo := &domains.ActionInfo{
				Message: &domains.Message{
					Role:    "user",
					Content: msg,
				},
				Result: &domains.ActionResult{
					Result: "",
					Type:   domains.ActionResultText,
				},
			}

			action := domains.NewMessageAction(chat)
			answer, err := action.Handle(ctx, actionInfo)
			if err != nil {
				logrus.Errorf("ChatGPT: chat error: %v", err)
			}

			logrus.Debugf("ChatGPT: chat: message=%s answer: %s", msg, answer.Result)

			result := escape.String(answer.Result)
			ctx.String(http.StatusOK, result)
		})

		logrus.Fatal(r.Run(port))*/
}
