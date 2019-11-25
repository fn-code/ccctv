package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotSetup struct {
	*tgbotapi.BotAPI
	GroupID int64
}

func New(apikey string, groupID int64) (*BotSetup, error) {
	bot, err := tgbotapi.NewBotAPI(apikey)
	if err != nil {
		return nil, err
	}
	return &BotSetup{bot, groupID}, nil
}

func (bot *BotSetup) SendMessage(text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(bot.GroupID, text)
	ms, err := bot.Send(msg)
	if err != nil {
		return tgbotapi.Message{}, err
	}
	return ms, nil
}
