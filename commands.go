package main

import (
	"database/sql"
	"fmt"
	"log"
	"weatherbot/pkg/dbsql"
	"weatherbot/pkg/notifications"
)

type Command interface {
	Execute() (string, error)
}

type CommandOut struct {
	profile *dbsql.UserProfile
	db      *sql.DB
}

func (c *CommandOut) Execute() (string, error) {
	c.profile.OutputFormat = swap(c.profile.OutputFormat)
	log.Printf("OutputFormat is %d", c.profile.OutputFormat)
	dbsql.SaveUserProfile(c.db, c.profile)
	return "Вы сменили формат вывода", nil
}

type CommandValue struct {
	profile *dbsql.UserProfile
	db      *sql.DB
}

func (c *CommandValue) Execute() (string, error) {
	c.profile.ValueFormat = swap(c.profile.ValueFormat)
	log.Printf("ValueFormat is %d", c.profile.ValueFormat)
	dbsql.SaveUserProfile(c.db, c.profile)
	return "Вы сменили размерность величин", nil
}

type CommandHelp struct{}

func (c *CommandHelp) Execute() (string, error) {
	temp, err := getText("./texts/help.txt")
	if err != nil {
		return "Я забыл, как я работаю", err
	}
	return temp, nil
}

type CommandStart struct{}

func (c *CommandStart) Execute() (string, error) {
	temp, err := getText("./texts/start.txt")
	if err != nil {
		return "Привет!", err
	}
	return temp, nil
}

type CommandDelNotif struct {
	userID int
}

func (c *CommandDelNotif) Execute() (string, error) {
	notifications.DeleteNotificationNOW(c.userID)
	return "Уведомления отключены", nil
}

type CommandSetNotif struct {
	request string
	userID  int
}

func (c *CommandSetNotif) Execute() (string, error) {
	City, Interval, err := checkNotificationComandFormat(c.request)
	if err != nil {
		return "", fmt.Errorf("неверный формат ввода")
	}
	notifications.SetNotification(c.userID, City, Interval)
	return "Уведомления подключены", nil
}

func SetCommand(input string, profile *dbsql.UserProfile, db *sql.DB, userID int) Command {
	command := sliceStringWithFirstSpace(input)
	switch command {
	case "/outputformat":
		return &CommandOut{
			profile: profile,
			db:      db,
		}
	case "/valueformat":
		return &CommandValue{
			profile: profile,
			db:      db,
		}
	case "/help":
		return &CommandHelp{}
	case "/start":
		return &CommandStart{}
	case "/deletenotification":
		return &CommandDelNotif{
			userID: userID,
		}
	default:
		// multiple names of command switch
		switch {
		case command == "/setnotification" || command == "/set":
			return &CommandSetNotif{
				request: input,
				userID:  userID,
			}
		default:
			return nil
		}
	}
}
