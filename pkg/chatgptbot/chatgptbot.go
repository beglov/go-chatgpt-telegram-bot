package chatgptbot

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Service struct {
	Bot             *tgbotapi.BotAPI
	client          *resty.Client
	telegramUserIds []int
}

func New(telegramBotToken string, openaiAPIKey string, telegramUserIds []int) (*Service, error) {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		return nil, err
	}

	client := resty.New()
	client.SetHeader("Authorization", "Bearer "+openaiAPIKey)

	s := Service{
		Bot:             bot,
		client:          client,
		telegramUserIds: telegramUserIds,
	}

	return &s, nil
}

func (s *Service) Start() {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := s.Bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		msg := s.Handler(update)
		if msg == nil {
			continue
		}
		if _, err := s.Bot.Send(msg); err != nil {
			log.Print("Error occurred while send msg: ", err)
			continue
		}
	}
}

// Handler выполняем обработку входящего сообщения.
// Возвращает сообщение которое требуется отправить или nil, если отправлять ничего не требуется.
func (s *Service) Handler(update tgbotapi.Update) (msg tgbotapi.Chattable) {
	if s.Bot.Debug {
		updateJSON, err := json.MarshalIndent(update, "", "  ")
		if err != nil {
			log.Printf("Error while marshal update: %s", err.Error())
		}
		log.Printf("Update received:\n%s\n", string(updateJSON))
	}

	switch {
	case update.Message != nil:
		msg = s.MessageHandler(update)
	default:
		log.Print("Unknown update type: ", update)
	}

	return msg
}

// MessageHandler выполняем обработку входящего сообщения.
// Возвращает сообщение которое требуется отправить или nil, если отправлять ничего не требуется.
func (s *Service) MessageHandler(update tgbotapi.Update) (msg tgbotapi.Chattable) {
	msg = s.Auth(update)
	if msg != nil {
		return msg
	}

	switch {
	case update.Message.IsCommand():
		msg = s.MessageCommandHandler(update)
	case update.Message.Text != "":
		msg = s.MessageTextHandler(update)
	default:
		log.Print("Unknown message type: ", update.Message)
	}

	return msg
}

// MessageCommandHandler выполняем обработку входящего сообщения.
// Возвращает сообщение которое требуется отправить или nil, если отправлять ничего не требуется.
func (s *Service) MessageCommandHandler(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		msg.Text = "Ask me anything you want"
	default:
		msg.Text = "Bot doesn't know this command"
	}

	return msg
}

// MessageTextHandler выполняем обработку входящего сообщения.
// Возвращает сообщение которое требуется отправить или nil, если отправлять ничего не требуется.
func (s *Service) MessageTextHandler(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	messages := []map[string]interface{}{
		{"role": "user", "content": update.Message.Text},
	}

	chatCompletion := &ChatCompletion{}

	// https://platform.openai.com/docs/api-reference/chat/create
	resp, err := s.client.R().
		SetBody(map[string]interface{}{
			"messages": messages,
			"model":    "gpt-3.5-turbo",
		}).
		SetResult(chatCompletion).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		log.Print(err)
		msg.Text = "Service Temporarily Unavailable"
		return msg
	}

	log.Print(resp)
	msg.Text = chatCompletion.Choices[0].Message.Content

	return msg
}

// Auth проверяет может ли пользователь взаимодействовать с ботом.
// Если переменная окружения TELEGRAM_USER_IDS не задана или пуста,
// то взаимодействовать с ботом разрешено всем пользователям.
// Иначе проверяем включен ли пользователь отправивший сообщение в список разрешенных.
// Возвращает nil если взаимодействие разрешено, иначе сообщение, которое будет отправлено пользователю.
func (s *Service) Auth(update tgbotapi.Update) tgbotapi.Chattable {
	if len(s.telegramUserIds) == 0 {
		return nil
	}

	for _, id := range s.telegramUserIds {
		if id == int(update.Message.From.ID) {
			return nil
		}
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied")
}
