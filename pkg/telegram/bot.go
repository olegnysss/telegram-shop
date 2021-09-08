package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegnysss/telebot_qiwi/pkg/couchbase"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"log"
)

type ID uint32

func InitBot(tgToken string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	return bot, updates
}

func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

//todo refactor input values
func HandleCommands(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, qiwiConfig qiwi.Config) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			handleStartCommand(update, bot)
		case "↪️Вернуться в главное меню ↩️":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Главное меню")
			msg.ReplyMarkup = StarterKeyboard
			sendMessage(bot, msg)
		case "Личный кабинет 👤":
			personalData := fmt.Sprintf(PersonalDataFormat, update.Message.From.ID, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, personalData)
			msg.ReplyMarkup = profileKeyboard
			sendMessage(bot, msg)
		case "💰 Пополнения":
			qiwi.CheckPayment(qiwiConfig)
		}
	}
}

func handleStartCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_, ok := couchbase.UsersMap[couchbase.ID(update.Message.Chat.ID)]
	if !ok {
		newUser := couchbase.User{
			UserId: couchbase.ID(update.Message.Chat.ID),
			Name:   update.Message.Chat.UserName,
		}
		couchbase.AppendUser(newUser)
	} else {
		log.Printf("Пользователь %d уже зарегистрирован", update.Message.Chat.ID)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать.")
	msg.ReplyMarkup = StarterKeyboard
	sendMessage(bot, msg)
}
