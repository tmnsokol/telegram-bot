package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/config"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/telegram/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string

	messages config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectUrl string, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectUrl, messages: messages}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()

	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)

			if err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}

			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}

	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.bot.GetUpdatesChan(u)

	if err != nil {
		return nil, err
	}

	return updates, nil

}
