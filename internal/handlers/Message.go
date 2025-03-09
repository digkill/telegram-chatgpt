package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gitlab.com/mediarise/appleclassbot/internal/commands"
	"gitlab.com/mediarise/appleclassbot/internal/components/chatGPT"
	"gitlab.com/mediarise/appleclassbot/internal/components/database"
	"gitlab.com/mediarise/appleclassbot/internal/components/redis"
	"gitlab.com/mediarise/appleclassbot/internal/config"
	"gitlab.com/mediarise/appleclassbot/internal/models"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const LIMIT_DAY_PROMPT int = 5

type MessageContext struct {
	Updater *UpdateTelegramData
	Payload string
	Config  *config.Config
}

type MessageHandler interface {
	Handle(message *tgbotapi.Message, ctx *MessageContext)
}

type InitHandler struct {
	Next MessageHandler
}

func (i *InitHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	if message.IsCommand() && message.Command() == "start" {
		err := ctx.Updater.SendMessageTelegram(
			message.Chat.ID,
			"Здарова, "+message.Chat.UserName+" гений в разработке! 😎\n\nЯ Архимед GPT, твой бро в мире математики! 🚀\n\n📌 Что умею?\n✅ Решаю любые задачи – алгебра, геометрия, уравнения, дроби, всё, что душа пожелает.\n✅ Фоткай задание – разберусь и объясню!\n✅ Помогу не только списать, но и реально понять, чтобы на контрольной ты был королём! 👑\n✅ Разжую даже самую жёсткую тему, как будто это мемчик с котиками.\n\n💬 Просто напиши мне вопрос или кинь фотку примера – и разберёмся на изи! 😏")

		if err != nil {
			logrus.Errorf("Cannot send message. Error: " + err.Error())
		}
	}

	i.Next.Handle(message, ctx)
}

type CommandMenuHandler struct {
	Next MessageHandler
}

func (h *CommandMenuHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	var databaseConfig = ctx.Config.DB
	var redisConfig = ctx.Config.Redis
	var db = database.NewDb(&databaseConfig)
	var user = models.NewUser(db)

	var limitDayPrompt = LIMIT_DAY_PROMPT

	userModel, _ := user.FindUserByUsername(message.Chat.UserName)

	stat, err := getStat(message, ctx, db)
	if err != nil {
		fmt.Println(err)
	}

	limitDayPrompt = limitDayPrompt + (stat * 10)

	if message.Command() == "start" {
		if userModel == nil {
			userModel, _ = user.CreateUser(message.Chat.UserName)
		}

		startCommand := commands.NewStartCommand(ctx.Updater.GetBot(), ctx.Config, db)
		err := startCommand.Execute(message, message.From.ID, message.Chat.UserName)
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	if message.Command() == "ref" {
		me, err := ctx.Updater.GetBot().GetMe()
		if err != nil {
			return
		}

		refCommand := commands.NewRefCommand(ctx.Updater.GetBot(), ctx.Config, db)
		refCommand.Execute(message.From.ID, me.UserName)

		return
	}

	if message.Command() == "stats" {
		stat, err := getStat(message, ctx, db)
		if err != nil {
			return
		}

		msg := tgbotapi.NewMessage(int64(message.From.ID), fmt.Sprintf("У вас %d рефералов! 🎉", stat))
		ctx.Updater.GetBot().Send(msg)

		return
	}

	if message.Command() == "report" {
		newReportCommand := commands.NewReportCommand(ctx.Updater.GetBot(), ctx.Config, db)
		st, _ := newReportCommand.Execute()

		// Вывод данных

		ctx.Updater.SendMessageTelegram(
			message.Chat.ID,
			"Статистика пользователей по дням:")

		for _, stat := range st {
			fmt.Printf("Дата: %s, Пользователь: %s, Запросов: %d, Всего команд: %d\n",
				stat.Date, stat.Username, stat.RequestCount, stat.TotalCount)

			ctx.Updater.SendMessageTelegram(
				message.Chat.ID,
				fmt.Sprintf("Дата: %s, Пользователь: %s, Запросов: %d, Всего команд: %d\n",
					stat.Date, stat.Username, stat.RequestCount, stat.TotalCount))

		}

		return

	}

	if message.Command() == "users" {
		newReportCommand := commands.NewUsersCommand(ctx.Updater.GetBot(), ctx.Config, db)
		st, _ := newReportCommand.Execute()

		// Вывод данных

		ctx.Updater.SendMessageTelegram(
			message.Chat.ID,
			"Статистика пользователей:")

		for _, stat := range st {
			fmt.Printf("Всего пользователей: %d\n",
				stat.TotalCount)

			ctx.Updater.SendMessageTelegram(
				message.Chat.ID,
				fmt.Sprintf("Всего пользователей: %d\n",
					stat.TotalCount))

		}

		return

	}

	if message.Command() == "author" {
		ctx.Updater.Handler.SendResult(
			message.Chat.ID,
			"<Vitaliy Edifanov> mediarise.ru",
			models.Button{},
		)
		return
	}

	if userModel != nil {
		var journalModel = models.NewJournal(database.NewDb(&databaseConfig))

		var newRedis = redis.NewRedis(&redisConfig)
		var keyUsername = message.Chat.UserName

		isHasData := newRedis.HasData(keyUsername)

		var count = 0
		if isHasData == false {
			err := newRedis.SetData(keyUsername, strconv.Itoa(0), time.Hour*24)
			if err != nil {
				log.Errorf("Ошибка установки значения: " + err.Error())
			}
		} else {
			count = int(newRedis.Increment(keyUsername))
		}

		var prompt = message.Text
		images := message.Photo

		if images != nil && len(*images) > 0 {
			photoId := (*images)[1].FileID

			fileId := tgbotapi.FileConfig{FileID: photoId}

			file, err := ctx.Updater.GetBot().GetFile(fileId)
			if err != nil {

				log.Errorf("Ошибка получения файла!: " + err.Error())
			}

			urlImage := file.Link(ctx.Updater.GetBot().Token)
			prompt = prompt + " " + urlImage

		}

		var journal, _ = journalModel.CreateJournal(userModel.Id, prompt, count)

		if journal != nil {

			if count > limitDayPrompt {

				err = ctx.Updater.SendMessageWithButtonsInRowToTelegram(
					message.Chat.ID,
					"Извините. Дневной лимит запросов исчерпан 😥",
					tgbotapi.NewInlineKeyboardButtonData("Получить запросы бесплатно 🤖", "ref_menu"),
				)

				if isHasData != false {
					newRedis.Decrement(keyUsername)
				}

				return
			}
		}

	}

	messObj, err := ctx.Updater.SendMessageTextTelegram(
		message.Chat.ID,
		"Решаю задачу 🤓...",
		tgbotapi.ModeMarkdown,
	)
	if err != nil {
		return
	}

	var chatGPTConfig = ctx.Config.ChatGPT

	chat := chatGPT.NewChatGPT(&chatGPTConfig)

	images := message.Photo

	var systemPrompt = "ответы пиши только на русском языке. Используй символы unicode, не используй нотацию LaTex, Tex. Начинаем новую тему, без учета предыдущих разговоров."

	if images != nil && len(*images) > 0 {
		photoId := (*images)[1].FileID

		fileId := tgbotapi.FileConfig{FileID: photoId}

		file, err := ctx.Updater.GetBot().GetFile(fileId)
		if err != nil {
			return
		}

		urlImage := file.Link(ctx.Updater.GetBot().Token)

		ext := filepath.Ext(urlImage)

		image, err := downloadFile(urlImage, ext)
		if err != nil {
			log.Fatal(err)
		}

		imgUrl := openai.ChatMessageImageURL{
			URL: image,
		}

		contentImg := openai.ChatMessagePart{
			ImageURL: &imgUrl,
			Type:     openai.ChatMessagePartTypeImageURL,
		}

		promptImage := "Реши задачу с картинки"
		if message.Text != "" {
			promptImage = message.Text
		}

		contentText := openai.ChatMessagePart{
			Text: promptImage,
			Type: openai.ChatMessagePartTypeText,
		}

		contentSystem := openai.ChatMessagePart{
			Text: systemPrompt,
			Type: openai.ChatMessagePartTypeText,
		}

		// Создаём JSON-объект в виде структуры
		data := []openai.ChatCompletionMessage{
			{
				Role:         "user",
				MultiContent: []openai.ChatMessagePart{contentImg, contentText},
			},
			{
				Role:         "system",
				MultiContent: []openai.ChatMessagePart{contentSystem},
			},
		}

		var contextGpt *gin.Context
		contextGpt = &gin.Context{}

		answer, err := chat.Chat(contextGpt, data)

		if err != nil {
			logrus.Error(err)
		}

		ctx.Updater.Handler.SendResult(
			message.Chat.ID,
			answer.Content,
			models.Button{
				Type: "show_main_menu",
			},
		)

	} else {

		contentText := openai.ChatMessagePart{
			Text: message.Text,
			Type: openai.ChatMessagePartTypeText,
		}

		contentSystem := openai.ChatMessagePart{
			Text: systemPrompt,
			Type: openai.ChatMessagePartTypeText,
		}

		// Создаём JSON-объект в виде структуры
		data := []openai.ChatCompletionMessage{
			{
				Role:         "user",
				MultiContent: []openai.ChatMessagePart{contentText},
			},
			{
				Role:         "system",
				MultiContent: []openai.ChatMessagePart{contentSystem},
			},
		}

		var contextGpt *gin.Context
		contextGpt = &gin.Context{}

		answer, err := chat.Chat(contextGpt, data)

		if err != nil {
			logrus.Error(err)
		}

		ctx.Updater.Handler.SendResult(
			message.Chat.ID,
			answer.Content,
			models.Button{
				Type: "show_main_menu",
			},
		)

	}
	_, err = ctx.Updater.Handler.RemoveMessage(message.Chat.ID, messObj.MessageID)
	if err != nil {
		log.Errorf("Ошибка удаления сообщения: " + err.Error())
	}

	return

	h.Next.Handle(message, ctx)
}

type FinishHandler struct{}

func (h *FinishHandler) Handle(message *tgbotapi.Message, ctx *MessageContext) {
	/*
		err := ctx.Updater.SendMessageTelegram(message.Chat.ID,
			"Проверьте корректность команды!")
		if err != nil {
			logrus.Errorf("Cannot send message. Error: " + err.Error())
		}

	*/
}

func downloadFile(url string, fileMimeType string) (string, error) {
	//Get the response bytes from the url
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return "", err
	}

	// Кодируем в Base64
	base64String, _ := EncodeImageToBase64(bodyBytes, fileMimeType)

	return base64String, nil
}

func EncodeImageToBase64(imageBytes []byte, fileMimeType string) (string, error) {

	// Кодируем в base64
	base64Str := base64.StdEncoding.EncodeToString(imageBytes)

	// Определяем MIME-тип по расширению
	mimeType := mime.TypeByExtension(fileMimeType)
	if mimeType == "" {
		mimeType = "application/octet-stream" // По умолчанию, если неизвестный тип
	}

	// Формируем data URL
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)

	return dataURL, nil
}

func getStat(message *tgbotapi.Message, ctx *MessageContext, db *database.DbComponent) (int, error) {
	statsCommand := commands.NewStatCommand(ctx.Updater.GetBot(), ctx.Config, db)
	count, err := statsCommand.Execute(message.From.ID)
	if err != nil {
		fmt.Println(err)
	}

	return count, nil
}
