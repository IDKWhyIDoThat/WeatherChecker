package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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

	lastUpdateCh := make(chan tgbotapi.Update)

	bot := createBot(botrefer)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	DB := dbsql.InitDB()
	defer DB.Close()

	go fetchLastUpdates(updates, lastUpdateCh)

	go checkNotificaions(DB, bot)

	for {
		lastUpdate := <-lastUpdateCh
		log.Printf("Update receiced")
		if lastUpdate.Message != nil {
			log.Printf("[%d] Author: %s Message: %s", lastUpdate.Message.Chat.ID, lastUpdate.Message.From.FirstName, lastUpdate.Message.Text)
			handleMessage(bot, lastUpdate, DB)
			lastUpdate = tgbotapi.Update{}
		}
		lastUpdateCh <- lastUpdate
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

func fetchLastUpdates(updates tgbotapi.UpdatesChannel, updatech chan<- tgbotapi.Update) {
	for update := range updates {
		updatech <- update
	}
}

func checkNotificaions(DB *sql.DB, bot *tgbotapi.BotAPI) {
	for {
		time.Sleep(notifyrefreshTime)
		ID, City, err := notifications.NotifyCheckout()
		if err == nil && ID != 0 {
			profile, err := dbsql.GetUserProfile(DB, ID)
			if err != nil {
				log.Print("error receiving data from server: ", err)
				continue
			}
			result, err := weather.GetCityWeatherData(City, profile.OutputFormat, profile.ValueFormat)
			if err != nil {
				log.Print("error receiving data from server: ", err)
			}
			sendMessageDirectly(bot, ID, result)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	bot.Send(msg)
}

func sendMessageDirectly(bot *tgbotapi.BotAPI, ID int, message string) {
	msg := tgbotapi.NewMessage(int64(ID), message)
	bot.Send(msg)
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, DB *sql.DB) {
	var response string
	if update.Message.Text[0] == '/' {
		log.Print("command detected")
		var err error
		response, err = handleCommand(update, DB)
		if err != nil {
			log.Print("command init failed: ", err)
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
		sendMessage(bot, update, response)
	}
}

func handleCommand(update tgbotapi.Update, db *sql.DB) (string, error) {
	userID := update.Message.From.ID
	profile, err := dbsql.GetUserProfile(db, userID)
	if err != nil {
		log.Fatal("Something gone terribly wrong")
	}

	request := update.Message.Text
	log.Printf("Request received: %s", request)

	switch {
	case request == "/outputformat":
		profile.OutputFormat = swap(profile.OutputFormat)
		log.Printf("OutputFormat is %d", profile.OutputFormat)
		dbsql.SaveUserProfile(db, profile)
		return "Вы сменили формат вывода", nil

	case request == "/valueformat":
		profile.ValueFormat = swap(profile.ValueFormat)
		log.Printf("ValueFormat is %d", profile.ValueFormat)
		dbsql.SaveUserProfile(db, profile)
		return "Вы сменили размерность величин", nil

	case request == "/help":
		temp, err := getText("./texts/help.txt")
		if err != nil {
			return "Я забыл, как я работаю", err
		}
		return temp, nil

	case request == "/start":
		temp, err := getText("./texts/start.txt")
		if err != nil {
			return "Привет!", err
		}
		return temp, nil

	case request == "/deletenotification":
		notifications.DeleteNotificationNOW(userID)
		return "Уведомления отключены", nil

	case strings.HasPrefix(request, "/setnotification") || strings.HasPrefix(request, "/set"):
		City, Interval, err := checkNotificationComandFormat(request)
		if err != nil {
			return "", fmt.Errorf("неверный формат ввода")
		}
		notifications.SetNotification(userID, City, Interval)
		return "Уведомления подключены", nil

	default:
		return "Неизвестная команда", fmt.Errorf("unknown command")
	}
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

func checkNotificationComandFormat(input string) (string, int, error) {
	parts := strings.Fields(input)
	if len(parts) != 3 {
		return "", 0, fmt.Errorf("неверный формат строки")
	}
	message := parts[1]
	duration, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", 0, fmt.Errorf("неверный формат числа")
	}
	return message, duration, nil
}
