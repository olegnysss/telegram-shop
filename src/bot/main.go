package main

import (
	"github.com/olegnysss/telebot_qiwi/pkg/config"
	"github.com/olegnysss/telebot_qiwi/pkg/couchbase"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"github.com/olegnysss/telebot_qiwi/pkg/telegram"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	//ToDo refactor to factory
	err = couchbase.ConnectToCouch(couchbaseConfig(cfg))
	if err != nil {
		log.Panic(err)
	}

	bot, updates := telegram.InitBot(cfg.TelegramToken)
	telegram.HandleCommands(updates, bot, qiwiConfig(cfg))
}

func couchbaseConfig(config *config.Config) couchbase.Config {
	return couchbase.Config{
		ConnString:    config.CouchConnString,
		CouchUsername: config.CouchUsername,
		CouchPassword: config.CouchPassword,
		BucketName:    config.CouchBucketName,
	}
}

func qiwiConfig(config *config.Config) qiwi.Config {
	return qiwi.Config{
		QiwiToken:        config.QiwiToken,
		QiwiWallet:       config.QiwiWallet,
		QiwiPaymentsPath: config.QiwiPaymentsPath,
	}
}
