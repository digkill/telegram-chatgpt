package telegram

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"sync"
)

type Telegram struct {
	*tgbotapi.BotAPI
}

type Wrapper struct {
	listOfTelegrams sync.Map
	debug           bool
}

func (wrapper *Wrapper) GetInstance(token string) (*Telegram, error) {
	t, ok := wrapper.listOfTelegrams.Load(token)
	if ok {
		return t.(*Telegram), nil
	}

	if !ok {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			return nil, err
		}

		bot.Debug = wrapper.debug

		logrus.Debugf("Authorized on account %s", bot.Self.UserName)
		wrapper.listOfTelegrams.Store(token, Telegram{bot})
	}

	t, ok = wrapper.listOfTelegrams.Load(token)
	if ok {
		return t.(*Telegram), nil
	}

	return nil, errors.New("error: could not find Telegram struct by token name")
}

func NewTelegram(isDebug bool) (*Wrapper, error) {

	return &Wrapper{
		debug: isDebug,
	}, nil
}
