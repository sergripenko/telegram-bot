package main

import (
	"telegram-bot/exchange_rates"
	"telegram-bot/models"
	"telegram-bot/services"
	"telegram-bot/services/db"
	"telegram-bot/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/gommon/log"
)

// кнопки над клавиатурой
var mainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		//tgbotapi.NewInlineKeyboardButtonURL("  Pornhub  ", "http://pornhub.com"),
		tgbotapi.NewInlineKeyboardButtonSwitch("Share bot", "open bot"),
		tgbotapi.NewInlineKeyboardButtonData("  Rates  ", "rates"),
		tgbotapi.NewInlineKeyboardButtonData("  Weather  ", "weather"),
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
	var err error
	var conf services.Config

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

	// TODO: send direct message
	//msg := tgbotapi.NewMessage(130421447,  "ffffff")
	//bot.Send(msg)
	//
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error(err)
	}

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	//time.Sleep(time.Second * 5)
	//updates.Clear()

	for update := range updates {

		// catch when user push the button
		if update.CallbackQuery != nil {
			//bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			//bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

			switch update.CallbackQuery.Data {
			case "weather":
				// check if user have saved location
				locExist := models.IfLocationExist(int(update.CallbackQuery.Message.Chat.ID))

				if !locExist {
					msg.Text = "Hi! Sorry, but I don't have your location :( \nSend me your location and I'll remember it"
					bot.Send(msg)

				} else {
					lat, lon := models.GetUsersCoords(int(update.CallbackQuery.Message.Chat.ID))
					weatherInfo, _ := weather.GetOpenWeather(int(update.CallbackQuery.Message.Chat.ID), lat, lon)
					msg.Text = weatherInfo
					log.Info("send weather")
					msg.ReplyMarkup = mainKeyboard
					bot.Send(msg)
				}

			case "rates":
				var rates string

				if rates, err = exchange_rates.GetRates(); err != nil {
					log.Error(err)
				}
				msg.Text = rates
				msg.ReplyMarkup = mainKeyboard
				log.Info("send rates")
				bot.Send(msg)
			}
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message != nil {
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "choose you category")

			user := &models.Users{
				UserName:  update.Message.Chat.UserName,
				FirstName: update.Message.Chat.FirstName,
				LastName:  update.Message.Chat.LastName,
				ChatID:    int(update.Message.Chat.ID),
			}
			// add new user
			models.AddNewUser(user)

			// case when user sent location
			if update.Message.Location != nil {
				log.Info("user send location")
				models.AddUsersLocation(int(update.Message.Location.Latitude), int(update.Message.Location.Longitude),
					int(update.Message.Chat.ID))

				weatherInfo, _ := weather.GetOpenWeather(int(update.Message.Chat.ID),
					int(update.Message.Location.Latitude), int(update.Message.Location.Longitude))
				msg.Text = weatherInfo
				log.Info("send weather info")
			}

			switch update.Message.Text {
			case "/start":
				msg.ReplyMarkup = mainKeyboard
			}
			msg.ReplyMarkup = mainKeyboard
			log.Info("send main keys")
			bot.Send(msg)
		}
		log.Info(update.Message.From.UserName, " sent ", update.Message.Text)

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
