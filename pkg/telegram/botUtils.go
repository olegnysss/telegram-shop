package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
)

var starterKeyboard = tgbotapi.NewReplyKeyboard(
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

var balanceKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Пополнить баланс 💰"),
		tgbotapi.NewKeyboardButton("Активировать купон"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("↪️Вернуться в главное меню ↩️"),
	),
)

var cashInKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Я пополнил баланс"),
		tgbotapi.NewKeyboardButton("🔙 Назад"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("↪️Вернуться в главное меню ↩️"),
	),
)

func getPaymentKeyboard(id int64, qiwiConfig qiwi.Config) tgbotapi.InlineKeyboardMarkup {
	cashInPath := fmt.Sprintf(qiwiConfig.QiwiCashInPath, qiwiConfig.QiwiWallet, id)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Перейти к оплате", cashInPath),
		),
	)
}

var personalDataFormat = "" +
	"➖➖➖➖➖➖➖➖➖➖\n" +
	"Ваш профиль:\n" +
	"🕶️ Ваш ID: %d\n" +
	"👏 Ваш никнейм: %s\n" +
	"➖➖➖➖➖➖➖➖➖➖"

var paymentInfo = "" +
	"➖➖➖➖➖➖➖➖➖➖\n" +
	"Информация об оплате:\n" +
	"🥝 QIWI-кошелек: +%s\n" +
	"📝 Комментарий к переводу: %d\n" +
	"➖➖➖➖➖➖➖➖➖➖\n" +
	"Внимание\n" +
	"Переводите ту сумму, на которую хотите пополнить баланс!\n" +
	"Заполняйте номер телефона и комментарий при переводе внимательно!\n" +
	"Администрация не несет ответственности за ошибочный перевод, возврата в данном случае не будет!\n" +
	"После перевода нажмите кнопку 'Я пополнил баланс'!\n" +
	"➖➖➖➖➖➖➖➖➖➖"
