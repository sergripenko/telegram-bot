package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/prometheus/common/log"
	"time"
)

// кнопки над клавиатурой
var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("  Pornhub  ", "http://pornhub.com"),
		tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

// кнопки вместо клавиатуры
//var numericKeyboard = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("1"),
//		tgbotapi.NewKeyboardButton("2"),
//		tgbotapi.NewKeyboardButton("3"),
//	),
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("4"),
//		tgbotapi.NewKeyboardButton("5"),
//		tgbotapi.NewKeyboardButton("6"),
//	),
//)

func main() {
	bot, err := tgbotapi.NewBotAPI("748950501:AAH3-jz_wZ4cmwNXZ6Ytd0kpXiazNn6tA0g")
	if err != nil {
		log.Error(err)
	}

	bot.Debug = false
	log.Info("Authorized on account ", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Error(err)
	}

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.CallbackQuery != nil {
			log.Info(update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			//switch update.Message.Text {
			//case "open":
			msg.ReplyMarkup = numericKeyboard

			//}
			bot.Send(msg)
		}

		log.Info(update.Message.From.UserName, " sent ", update.Message.Text)

		// handling user commands
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /sayhi or /status."
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)

			//} else {
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "you say "+ update.Message.Text)
			//	msg.ReplyToMessageID = update.Message.MessageID
			//
			//	if _, err = bot.Send(msg); err != nil {
			//		log.Error(err)
			//	}
		}
	}
}
