package main

import (
	"telegram-bot/models"
	"telegram-bot/services"
	"telegram-bot/services/db"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/gommon/log"
)

// кнопки над клавиатурой
var mainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

func main() {
	// init DB
	db.InitDB()

	var conf services.Config
	var err error
	if conf, err = services.GetConfig(); err != nil {
		log.Error(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
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
		//if update.CallbackQuery != nil {
		//	log.Info(update)
		//
		//	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
		//
		//	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		//
		//}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			log.Info(update.Message.Chat.LastName)
			log.Info(update.Message.Chat.FirstName)
			log.Info(update.Message.Chat.ID)

			user := &models.Users{
				UserName:  update.Message.Chat.UserName,
				FirstName: update.Message.Chat.FirstName,
				LastName:  update.Message.Chat.LastName,
				ChatID:    int(update.Message.Chat.ID),
			}
			fl, err := models.AddUser(user)
			log.Error(err)
			log.Error(fl)
			//dbcon.NewRecord(user)// => returns `true` as primary key is blank
			//dbcon.Create(&user)
			//dbcon.NewRecord(user) // => return `false` after `user` created

			switch update.Message.Text {
			case "/start":
				msg.ReplyMarkup = mainKeyboard
				log.Info(44)
				bot.Send(msg)
			}
		}

		log.Info(update.Message.From.UserName, " sent ", update.Message.Text)

		//handling user commands
		//if update.Message.IsCommand() {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//	switch update.Message.Command() {
		//	case "help":
		//		msg.Text = "type /sayhi or /status."
		//	//case "sayhi":
		//	//	msg.Text = "Hi :)"
		//	case "start":
		//		msg.ReplyMarkup = mainKeyboard
		//		bot.Send(msg)
		//		log.Info(44)
		//	default:
		//		msg.Text = "I don't know that command"
		//	}
		//	bot.Send(msg)

		//} else {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "you say "+ update.Message.Text)
		//	msg.ReplyToMessageID = update.Message.MessageID
		//
		//	if _, err = bot.Send(msg); err != nil {
		//		log.Error(err)
		//	}
		//}
		//}
	}

}
