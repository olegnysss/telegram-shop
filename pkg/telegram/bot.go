package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegnysss/telebot_qiwi/pkg/couchbase"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"log"
	"os"
	"strconv"
)

type ID uint32

var Logger *log.Logger

type Bot struct {
	bot   *tgbotapi.BotAPI
	couch *couchbase.CouchClient
	qiwi  *qiwi.QiwiClient
}

func NewBot(bot *tgbotapi.BotAPI, couch *couchbase.CouchClient, qiwi *qiwi.QiwiClient) *Bot {
	return &Bot{
		bot:   bot,
		couch: couch,
		qiwi:  qiwi,
	}
}

func (b *Bot) Start() {
	err := initLogs()
	if err != nil {
		log.Println(err)
	}

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	b.HandleCommands(updates)
}

func (b *Bot) HandleCommands(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		Logger.Printf("Message: %+v\nFrom Chat %+v", update.Message, update.Message.Chat)

		chatId := update.Message.Chat.ID
		userName := update.Message.Chat.UserName

		switch update.Message.Text {
		case "/start":
			_, err := b.couch.UsersAdapter.CheckUser(chatId, userName)
			if err != nil {
				log.Println(err)
			}
			msg := tgbotapi.NewMessage(chatId, "Добро пожаловать.")
			msg.ReplyMarkup = starterKeyboard
			b.sendMessage(msg)
		case "↪️Вернуться в главное меню ↩️":
			msg := tgbotapi.NewMessage(chatId, "Главное меню")
			msg.ReplyMarkup = starterKeyboard
			b.sendMessage(msg)
		case "Личный кабинет 👤":
			personalData := fmt.Sprintf(personalDataFormat, chatId, userName)
			msg := tgbotapi.NewMessage(chatId, personalData)
			msg.ReplyMarkup = profileKeyboard
			b.sendMessage(msg)
		case "💰 Пополнения":
			transactions, err := b.couch.TransactionsAdapter.FetchTransactions(chatId)
			if err != nil {
				log.Println(err)
			}
			msgText := fmt.Sprintf("Список ваших пополнений:\n")
			i := 1
			for _, transaction := range transactions {
				msgText += fmt.Sprintf("%d. Пополнение баланса через киви на сумму: %.0f \n", i, transaction.Sum)
				i++
			}
			msg := tgbotapi.NewMessage(chatId, msgText)
			b.sendMessage(msg)
		case "Баланс 💰":
			user, err := b.couch.UsersAdapter.CheckUser(chatId, userName)
			if err != nil {
				log.Println(err)
			}
			message := fmt.Sprintf("Ваш баланс: %.0f", user.Balance)
			msg := tgbotapi.NewMessage(chatId, message)
			msg.ReplyMarkup = balanceKeyboard
			b.sendMessage(msg)
		case "Пополнить баланс 💰":
			message := fmt.Sprintf(paymentInfo, b.qiwi.Wallet, chatId)
			msg := tgbotapi.NewMessage(chatId, message)
			msg.ReplyMarkup = cashInKeyboard
			b.sendMessage(msg)
			msg1 := tgbotapi.NewMessage(chatId, "Для пополнения баланса нажмите кнопку ниже")
			msg1.ReplyMarkup = getPaymentKeyboard(chatId, b.qiwi)
			b.sendMessage(msg1)
		case "Я пополнил баланс":
			// timeout 60 sec
			payResp, err := b.qiwi.CheckPayment()
			if err != nil {
				log.Println(err)
			}
			transactions, err := b.couch.TransactionsAdapter.ParseTransactions(strconv.FormatInt(chatId, 10), payResp)
			isNew, err := b.couch.TransactionsAdapter.ProcessTnx(transactions, b.couch)
			msg := tgbotapi.NewMessage(chatId, "")
			if err != nil {
				couchbase.Logger.Println(err)
			}
			if isNew {
				couchbase.Logger.Printf("Transactions stored in database. %+v", transactions)
				msg.Text = "Пополнено!"
			} else {
				msg.Text = "Новых платежей не обнаружено"
			}
			b.sendMessage(msg)
		}
	}
}

func (b *Bot) sendMessage(msg tgbotapi.MessageConfig) {
	if _, err := b.bot.Send(msg); err != nil {
		log.Println(err)
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
