package domains

import (
	"context"
	"fmt"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/pkg/escape"
	"github.com/sirupsen/logrus"
)

var _ Action = &MessageAction{}

type MessageAction struct {
	Chat Chat
}

func NewMessageAction(chat Chat) *MessageAction {
	return &MessageAction{Chat: chat}
}

func (a MessageAction) Handle(ctx context.Context, actionInfo *ActionInfo) (res *ActionResult, err error) {

	msg := actionInfo.GetText()
	messages := MakeMessages(msg)

	fmt.Println(messages)

	answer, err := a.Chat.Chat(ctx, messages)
	if err != nil {
		logrus.Errorf("ChatGPT: chat error: %v", err)
	}

	logrus.Debugf("ChatGPT: chat: message=%s answer: %s", msg, answer.Content)

	result := escape.String(answer.Content)

	actionInfo.Result = &ActionResult{
		Result: result,
		Type:   ActionResultText,
	}
	return actionInfo.Result, nil
}

func MakeMessages(content string) []Message {

	return []Message{
		{
			Role:    "user",
			Content: content,
		},
	}
}
