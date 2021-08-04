package main

import (
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/config"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/server"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/telegram"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/telegram/repository"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/telegram/repository/boltdb"
	pocket "github.com/zhashkevych/go-pocket-sdk"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Panic(err)
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)

	if err != nil {
		log.Fatal(err)
	}

	db, err := initDb(cfg)

	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerUrl, cfg.Messages)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotUrl)

	go func() {
		if err := authorizationServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDb(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
