package main

import (
	"github.com/beglov/go-chatgpt-telegram-bot/pkg/chatgptbot"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type server struct {
	bot *chatgptbot.Service
}

func main() {
	srv := new()
	srv.run()
}

// new создаёт объект и службы сервера и возвращает указатель на него.
func new() *server {
	godotenv.Load(".env") //nolint:errcheck

	telegramBotToken := getTelegramBotToken()
	openaiAPIKey := getOpenaiAPIKey()
	telegramUserIds := getTelegramUserIds()
	retentionPeriod := getMessagesRetentionPeriod()

	bot, err := chatgptbot.New(telegramBotToken, openaiAPIKey, telegramUserIds, retentionPeriod)
	if err != nil {
		log.Fatal(err)
	}

	bot.Bot.Debug = true

	log.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	return &server{
		bot: bot,
	}
}

func getTelegramBotToken() string {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is blank")
	}
	return telegramBotToken
}

func getOpenaiAPIKey() string {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		log.Fatal("OPENAI_API_KEY is blank")
	}
	return openaiAPIKey
}

func getTelegramUserIds() []int {
	telegramUserIds := os.Getenv("TELEGRAM_USER_IDS")
	if telegramUserIds == "" {
		return make([]int, 0)
	}

	// Split string by comma separator
	strSlice := strings.Split(telegramUserIds, ",")
	// Convert each string element to integer and append to slice
	intSlice := make([]int, len(strSlice))
	for i, s := range strSlice {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal("Failed to convert", s, " to integer")
		}
		intSlice[i] = n
	}

	return intSlice
}

func getMessagesRetentionPeriod() int {
	retentionPeriod := os.Getenv("MESSAGES_RETENTION_PERIOD")
	if retentionPeriod == "" {
		retentionPeriod = "15"
	}

	i, err := strconv.Atoi(retentionPeriod)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

// run запускает бота.
func (srv *server) run() {
	srv.bot.Start()
}
