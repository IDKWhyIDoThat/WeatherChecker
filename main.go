package main

import (
	"log"
	"os"
	"time"
	"weatherbot/pkg/additional/getsmth"
	"weatherbot/pkg/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
	u.Timeout = 60

	var lastUpdate tgbotapi.Update

	updates, _ := bot.GetUpdatesChan(u)
	go func() {
		for update := range updates {
			lastUpdate = update
		}
	}()

	for {
		if lastUpdate.Message != nil {
			log.Printf("[%d] Author: %s Message: %s", lastUpdate.Message.Chat.ID, lastUpdate.Message.From.FirstName, lastUpdate.Message.Text)
			handleMessage(bot, lastUpdate)
			lastUpdate = tgbotapi.Update{} // Сбросить lastUpdate после обработки сообщения
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	city := ""
	outputformat := 2
	valueformat := 1 // 1 is for KM, 2 is for M
	city = update.Message.Text
	result, err := weather.GetCityWeatherData(city, outputformat, valueformat)
	if err != nil {
		log.Print("error receiving data from server: ", err)
	}
	if result != "" {
		sendMessage(bot, update, result)
	}
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	time.Sleep(500 * time.Millisecond)
	bot.Send(msg)
}
