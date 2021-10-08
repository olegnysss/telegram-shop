package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegnysss/telebot_qiwi/pkg/config"
	"github.com/olegnysss/telebot_qiwi/pkg/couchbase"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"github.com/olegnysss/telebot_qiwi/pkg/telegram"
	"log"
	"os"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Println(err)
	}
	logsInit()

	botApi, err := tgbotapi.NewBotAPI(cfg.Telegram.TelegramToken)
	if err != nil {
		log.Println(err)
	}

	couch := couchbase.InitCouchClient(cfg.Couch)
	_, err = couch.ConnectToCouch()
	if err != nil {
		log.Println(err)
	}

	qiwiClient := qiwi.InitQiwiClient(cfg.Qiwi)
	newBot := telegram.NewBot(botApi, couch, qiwiClient)
	newBot.Start()
}

func logsInit() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}
