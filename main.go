package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"weatherbot/pkg/additional/getsmth"
	"weatherbot/pkg/additional/support"
	"weatherbot/pkg/dbsql"
	"weatherbot/pkg/notifications"
	"weatherbot/pkg/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const refreshTime = 300 * time.Millisecond
const notifyrefreshTime = time.Minute
const logfilename = "bot.log"
const botrefer = "t.me/IDKWChekcer_bot"

func main() {
	logfile, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	bot := createBot(botrefer)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	DB := dbsql.InitDB()
	defer DB.Close()

	var lastUpdate tgbotapi.Update

	go checkNotificaions(DB, bot)

	go func() {
		for update := range updates {
			lastUpdate = update
		}
	}()

	for {
		if lastUpdate.Message != nil {
			log.Printf("Update receiced")
			log.Printf("[%d] Author: %s Message: %s", lastUpdate.Message.Chat.ID, lastUpdate.Message.From.FirstName, lastUpdate.Message.Text)
			handleMessage(bot, lastUpdate, DB)
			lastUpdate = tgbotapi.Update{}
		}
		time.Sleep(refreshTime)
	}
}

func createBot(refer string) *tgbotapi.BotAPI {
	BOTToken, err := getsmth.GetAPIkey(refer)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := tgbotapi.NewBotAPI(BOTToken)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func checkNotificaions(DB *sql.DB, bot *tgbotapi.BotAPI) {
	for {
		time.Sleep(notifyrefreshTime)
		ID, City, err := notifications.NotifyCheckout()
		if err == nil && ID != 0 {
			profile, err := dbsql.GetUserProfile(DB, ID)
			if err != nil {
				log.Print("error receiving data from DB: ", err)
				continue
			}
			result, err := weather.GetCityWeatherData(City, profile.OutputFormat, profile.ValueFormat)
			if err != nil {
				log.Print("error receiving data from server: ", err)
			}
			sendMessage(bot, int64(ID), result)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, ID int64, message string) {
	msg := tgbotapi.NewMessage(ID, message)
	bot.Send(msg)
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, DB *sql.DB) {
	var response string
	if update.Message.Text[0] == '/' {
		log.Print("command detected")
		var err error
		response, err = handleCommand(update, DB)
		if err != nil {
			log.Printf("command exec or init error: %s", err)
		}
	} else {
		city := update.Message.Text
		city = strings.ReplaceAll(city, " ", "-")

		profile, err := dbsql.GetUserProfile(DB, int(update.Message.Chat.ID))
		if err != nil {
			log.Print("error receiving data from DB: ", err)
			return
		}

		if support.HasCyrillic(city) {
			response = "Название города должно быть написано на английском"
		} else {
			response, err = weather.GetCityWeatherData(city, profile.OutputFormat, profile.ValueFormat)
			if err != nil {
				log.Print("error receiving data from server: ", err)
			}
		}
	}
	if response != "" {
		sendMessage(bot, update.Message.Chat.ID, response)
	}
}

func handleCommand(update tgbotapi.Update, DB *sql.DB) (string, error) {
	profile, err := dbsql.GetUserProfile(DB, update.Message.From.ID)
	if err != nil {
		log.Fatal("Something gone terribly wrong")
	}
	MyCommand := SetCommand(update.Message.Text, profile, DB, update.Message.From.ID)
	if MyCommand == nil {
		return "Неизвестная команда", fmt.Errorf("unknown command")
	}
	response, err := MyCommand.Execute()
	if err != nil {
		return "Произошла ошибка в ходе выполнения команды", fmt.Errorf("command error: %s", err)
	}
	return response, nil
}
