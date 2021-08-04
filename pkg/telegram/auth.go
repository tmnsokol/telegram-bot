package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/youtube-telegram-bot/youtube-telegram-bot/pkg/telegram/repository"
)

func (bot *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := bot.generateAuthorizationLink(message.Chat.ID)

	if err != nil {
		return err
	}

	msgText := fmt.Sprintf(bot.messages.Responses.Start, authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err = bot.bot.Send(msg)
	return err
}
func (bot *Bot) getAccessToken(chatID int64) (string, error) {
	return bot.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (bot *Bot) generateAuthorizationLink(chatID int64) (string, error) {

	redirectURL := bot.generateRedirectUrl(chatID)

	requestToken, err := bot.pocketClient.GetRequestToken(context.Background(), redirectURL)

	if err != nil {
		return "", err
	}

	if err := bot.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return bot.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectUrl(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
