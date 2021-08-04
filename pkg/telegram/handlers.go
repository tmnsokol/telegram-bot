package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart = "start"
)

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)

	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	access_token, err := b.getAccessToken(message.Chat.ID)

	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	_, err = url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: access_token,
		URL:         message.Text,
	}); err != nil {
		return errInvalidURL
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccesfully)
	b.bot.Send(msg)
	return err
}
