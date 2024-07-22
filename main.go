package main

import (
	"database/sql"
	"fmt"
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
	u.Timeout = 300

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
	if update.Message.Text[0] == '/' {
		log.Print("command detected")
		err := handleCommand(bot, update, DB)
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

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) error {
	userID := update.Message.From.ID
	profile, err := dbsql.GetUserProfile(db, userID)
	if err != nil {
		log.Fatal("Something gone terribly wrong")
	}

	request := update.Message.Text
	log.Printf("Request received: %s", request)
	switch request {
	case "/outputformat":
		profile.OutputFormat = swap(profile.OutputFormat)
		log.Printf("OutputFormat is %d", profile.OutputFormat)
		dbsql.SaveUserProfile(db, profile)
	case "/valueformat":
		profile.ValueFormat = swap(profile.ValueFormat)
		log.Printf("ValueFormat is %d", profile.ValueFormat)
		dbsql.SaveUserProfile(db, profile)
	case "/start":
		dbsql.SaveUserProfile(db, profile)
		temp, err := getText("./texts/start.txt")
		if err != nil {
			return err
		}
		sendMessage(bot, update, temp)
	default:
		return fmt.Errorf("unknown command")
	}
	return nil
}

func getText(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func swap(x int) int {
	if x == 1 {
		return 2
	} else {
		return 1
	}
}
