package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegnysss/telebot_qiwi/pkg/couchbase"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"log"
	"os"
)

type ID uint32

var Logger *log.Logger

func InitBot(tgToken string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.OpenFile("telegram.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(file, "TG: ", log.Ldate|log.Ltime|log.Lshortfile)

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

		Logger.Printf("Message: %+v\nFrom Chat %+v", update.Message, update.Message.Chat)

		switch update.Message.Text {
		case "/start":
			_, err := checkUser(update)
			if err != nil {
				log.Panic(err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать.")
			msg.ReplyMarkup = StarterKeyboard
			sendMessage(bot, msg)
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
		case "Баланс 💰":
			user, err := checkUser(update)
			if err != nil {
				log.Panic(err)
			}
			message := fmt.Sprintf("Ваш баланс: %f", user.Balance)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			sendMessage(bot, msg)
		}
	}
}

func checkUser(update tgbotapi.Update) (couchbase.User, error) {
	user, ok := couchbase.UsersMap[couchbase.ID(update.Message.Chat.ID)]
	if !ok {
		newUser := couchbase.User{
			UserId: couchbase.ID(update.Message.Chat.ID),
			Name:   update.Message.Chat.UserName,
		}
		return couchbase.AppendUser(newUser)
	} else {
		log.Printf("Пользователь %d уже зарегистрирован", update.Message.Chat.ID)
		return user, nil
	}
}
