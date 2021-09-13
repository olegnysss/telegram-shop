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

	err = initLogs()
	if err != nil {
		log.Panic(err)
	}

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

// HandleCommands todo refactor input values
func HandleCommands(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, qiwiConfig qiwi.Config) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		Logger.Printf("Message: %+v\nFrom Chat %+v", update.Message, update.Message.Chat)

		chatId := update.Message.Chat.ID
		userName := update.Message.Chat.UserName

		switch update.Message.Text {
		case "/start":
			_, err := couchbase.CheckUser(chatId, userName)
			if err != nil {
				log.Panic(err)
			}
			msg := tgbotapi.NewMessage(chatId, "Добро пожаловать.")
			msg.ReplyMarkup = starterKeyboard
			sendMessage(bot, msg)
		case "↪️Вернуться в главное меню ↩️":
			msg := tgbotapi.NewMessage(chatId, "Главное меню")
			msg.ReplyMarkup = starterKeyboard
			sendMessage(bot, msg)
		case "Личный кабинет 👤":
			personalData := fmt.Sprintf(personalDataFormat, update.Message.From.ID, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(chatId, personalData)
			msg.ReplyMarkup = profileKeyboard
			sendMessage(bot, msg)
		case "💰 Пополнения":
			qiwi.CheckPayment(qiwiConfig)
		case "Баланс 💰":
			user, err := couchbase.CheckUser(chatId, userName)
			if err != nil {
				log.Panic(err)
			}
			message := fmt.Sprintf("Ваш баланс: %f", user.Balance)
			msg := tgbotapi.NewMessage(chatId, message)
			msg.ReplyMarkup = balanceKeyboard
			sendMessage(bot, msg)
		case "Пополнить баланс 💰":
			_, err := couchbase.FetchTransactions(chatId)
			if err != nil {
				log.Panic(err)
			}
			message := fmt.Sprintf(paymentInfo, qiwiConfig.QiwiWallet, chatId)
			msg := tgbotapi.NewMessage(chatId, message)
			msg.ReplyMarkup = cashInKeyboard
			sendMessage(bot, msg)
			msg1 := tgbotapi.NewMessage(chatId, "Для пополнения баланса нажмите кнопку ниже")
			msg1.ReplyMarkup = getPaymentKeyboard(chatId, qiwiConfig)
			sendMessage(bot, msg1)
		case "Я пополнил баланс":

		}
	}
}

func initLogs() error {
	file, err := os.OpenFile("telegram.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	Logger = log.New(file, "TG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return err
}
