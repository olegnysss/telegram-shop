package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api"

var StarterKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Купить 💰"),
		tgbotapi.NewKeyboardButton("Наличие товаров"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Личный кабинет 👤"),
		tgbotapi.NewKeyboardButton("Баланс 💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Помощь 🤔"),
		tgbotapi.NewKeyboardButton("Правила 📜"),
		tgbotapi.NewKeyboardButton("О боте 🔒"),
	),
)

var profileKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💰 Пополнения"),
		tgbotapi.NewKeyboardButton("🛒 Покупки"),
		tgbotapi.NewKeyboardButton("👤 Рефералы"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("↪️Вернуться в главное меню ↩️"),
	),
)

var PersonalDataFormat = "" +
	"➖➖➖➖➖➖➖➖➖➖\n" +
	"Ваш профиль:\n" +
	"🕶️ Ваш ID: %d\n" +
	"👏 Ваш никнейм: %s\n" +
	"➖➖➖➖➖➖➖➖➖➖"

