package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	"weatherbot/pkg/additional/getsmth"
	"weatherbot/pkg/dbsql"
	"weatherbot/pkg/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const refreshTime = 500 * time.Millisecond

func main() {
	file, err := os.OpenFile("bot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	BOTToken, err := getsmth.GetAPIkey("t.me/IDKWChekcer_bot")

	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(BOTToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	var lastUpdate tgbotapi.Update

	updates, _ := bot.GetUpdatesChan(u)
	go func() {
		for update := range updates {
			lastUpdate = update
		}
	}()

	DB := dbsql.InitDB()

	for {
		if lastUpdate.Message != nil {
			log.Printf("[%d] Author: %s Message: %s", lastUpdate.Message.Chat.ID, lastUpdate.Message.From.FirstName, lastUpdate.Message.Text)
			handleMessage(bot, lastUpdate, DB)
			lastUpdate = tgbotapi.Update{} // cбросить lastUpdate после обработки сообщения
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, DB *sql.DB) {
	if string(update.Message.Text[0]) == `\` {
		err := dbsql.HandleCommand(DB, update)
		if err != nil {
			log.Print("command init failed: ", err)
		}
		return
	}
	city := update.Message.Text
	profile, err := dbsql.GetUserProfile(DB, int(update.Message.Chat.ID))
	if err != nil {
		log.Print("error receiving data from DB: ", err)
		return
	}
	result, err := weather.GetCityWeatherData(city, profile.OutputFormat, profile.ValueFormat)
	if err != nil {
		log.Print("error receiving data from server: ", err)
	}
	if result != "" {
		sendMessage(bot, update, result)
	}
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	time.Sleep(refreshTime)
	bot.Send(msg)
}
