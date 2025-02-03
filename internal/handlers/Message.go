package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/digkill/telegram-chatgpt/internal/models"
	"github.com/digkill/telegram-chatgpt/internal/services/chatgpt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"mime"
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

			//	ctxs := context.Background()

			//	bytes, err := openaiClient.CreateFileBytes(ctxs, openai.FileBytesRequest{
			//		Name:    "Пример",
			//		Bytes:   []byte(image),
			//		Purpose: openai.PurposeFineTune,
			//	})

			if err != nil {
				return
			}

			fmt.Println("🤓🤓🤓🤓🤓")
			//	fmt.Println(bytes)
			fmt.Println("🤓🤓🤓🤓🤓")

			//chat.GetChat().CreateFile()

			// ChatMessageImageURL

			// Определяем структуры для JSON

			imgUrl := openai.ChatMessageImageURL{
				URL: image,
			}

			contentImg := openai.ChatMessagePart{
				ImageURL: &imgUrl,
				Type:     openai.ChatMessagePartTypeImageURL,
			}

			contentText := openai.ChatMessagePart{
				Text:     "Реши задачу и результат приведи в формат HTML допустимы только теги <b>, <i>, <u>, <s>, <span>, <a>, <pre>, <code>, <blockquote>, другие HTML Теги использовать запрещено, приведи формулы в удобном виде с LaTeX, чтобы они были хорошо читаемыми,",
				ImageURL: &imgUrl,
				Type:     openai.ChatMessagePartTypeText,
			}

			// Создаём JSON-объект в виде структуры
			data := []openai.ChatCompletionMessage{
				{
					Role:         "user",
					MultiContent: []openai.ChatMessagePart{contentImg, contentText},
				},
			}

			//jsonData, err := json.Marshal(data)

			//	fmt.Println(string(jsonData))

			if err != nil {
				fmt.Println("Ошибка при кодировании в JSON:", err)
				return
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
			return

		} else {
			fmt.Println("Нет доступных изображений")
		}

		//	actionInfo := domains.ActionInfo{
		//		Message: &domains.Message{Role: "User", Content: mm},
		//	}

		// msg := actionInfo.GetText()
		//	messages := domains.MakeMessages(msg)

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
	base64String, _ := EncodeImageToBase64(bodyBytes)

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

func EncodeImageToBase64(imageBytes []byte) (string, error) {

	// Кодируем в base64
	base64Str := base64.StdEncoding.EncodeToString(imageBytes)

	// Определяем MIME-тип по расширению
	//ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(".jpg")
	if mimeType == "" {
		mimeType = "application/octet-stream" // По умолчанию, если неизвестный тип
	}

	// Формируем data URL
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)

	return dataURL, nil
}
