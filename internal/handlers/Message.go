package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/digkill/telegram-chatgpt/internal/domains"
	"github.com/digkill/telegram-chatgpt/internal/models"
	"github.com/digkill/telegram-chatgpt/internal/services/chatgpt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
)

type MessageContext struct {
	Updater *UpdateTelegramData
	Payload string
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
	if message.Command() == "start" {

		/*ctx.Updater.Handler.SendListMenu(
			message.Chat.ID,
			"Выберите услугу",
			models.Button{
				Type: "show_main_menu",
			},
		)
		*/

		return
	} else {

		err := ctx.Updater.SendMessageTelegram(
			message.Chat.ID,
			"Решаю задачу 🤓...",
		)
		if err != nil {
			return
		}

		var openAIToken = os.Getenv("CHATGPT_TOKEN")
		var openAIURL = os.Getenv("CHATGPT_URL")

		config := openai.DefaultConfig(openAIToken)
		if openAIURL != "" {
			config.BaseURL = openAIURL
		}

		openaiClient := openai.NewClientWithConfig(config)
		chat := chatgpt.NewChatGPT(openaiClient)

		images := message.Photo

		if images != nil && len(*images) > 0 {
			photoId := (*images)[0].FileID

			fileId := tgbotapi.FileConfig{FileID: photoId}

			file, err := ctx.Updater.GetBot().GetFile(fileId)
			if err != nil {
				return
			}

			urlImage := file.Link(ctx.Updater.GetBot().Token)

			image, err := downloadFile(urlImage)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("🤓🤓🤓🤓🤓")
			fmt.Println(image)
			fmt.Println("🤓🤓🤓🤓🤓")

			//	ctxs := context.Background()

			//	bytes, err := openaiClient.CreateFileBytes(ctxs, openai.FileBytesRequest{
			//		Name:    "Пример",
			//		Bytes:   []byte(image),
			//		Purpose: openai.PurposeFineTune,
			//	})

			fmt.Println("🤓🤓🤓🤓🤓")
			fmt.Println(err)
			fmt.Println("🤓🤓🤓🤓🤓")

			if err != nil {
				return
			}

			fmt.Println("🤓🤓🤓🤓🤓")
			//	fmt.Println(bytes)
			fmt.Println("🤓🤓🤓🤓🤓")

			//chat.GetChat().CreateFile()

			// ChatMessageImageURL

			message.Text = "Реши задачу по ссылке картинки " + urlImage
		} else {
			fmt.Println("Нет доступных изображений")
		}

		actionInfo := domains.ActionInfo{
			Message: &domains.Message{Role: "User", Content: message}
		}

		msg := actionInfo.GetText()
		messages := domains.MakeMessages(msg)

		var contextGpt *gin.Context
		contextGpt = &gin.Context{}

		answer, err := chat.Chat(contextGpt, messages)

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
		return
	}

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

func downloadFile(url string) (string, error) {
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
	base64String := base64.StdEncoding.EncodeToString(bodyBytes)

	//bodyBytes = []byte(strings.ToValidUTF8(string(bodyBytes), ""))
	//fmt.Println("!!!!")
	//fmt.Println(string(bodyBytes))
	//if response.StatusCode != 200 {
	//	return errors.New("Received non 200 response code")
	//	}
	//Create a empty file

	//Write the bytes to the fiel
	//_, err = io.Copy(file, response.Body)

	// Парсим JSON как массив
	// var records []map[string]interface{}
	//	if err := json.Unmarshal(bodyBytes, &records); err != nil {
	//	//	fmt.Println("Ошибка JSON:", err)
	//	}

	//	var jsonlData string

	// Записываем каждую строку как отдельный JSON-объект
	//	for _, record := range records {
	//	line, _ := json.Marshal(record) // Конвертируем в JSON строку
	//		jsonlData += string(line) + "\n"
	//}

	//fmt.Println("Файл успешно конвертирован в .jsonl!")
	//	fmt.Println(jsonlData)

	return base64String, nil
}
