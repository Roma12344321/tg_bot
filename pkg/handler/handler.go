package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"myBot/pkg/service"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitBot(ctx context.Context) {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("token"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)
	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			if update.Message != nil {
				if update.Message.IsCommand() {
					h.handleCommand(bot, update.Message)
				} else {
					h.handleMessage(bot, update.Message)
				}
			}
		}
	}
}

func (h *Handler) handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "login":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите логин и пароль в формате <ЛОГИН> <ПАРОЛЬ>")
		bot.Send(msg)
	}
}

func (h *Handler) handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	str := strings.Split(message.Text, " ")
	if len(str) < 2 || len(str) > 2 {
		h.handleStandardMessage(bot, message)
		return
	}
	err := h.service.AuthService.LogIn(message.Chat.ID, str[0], str[1])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "успешно")
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func (h *Handler) handleStandardMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вот список доступных команд:\n/login")
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
