package dbsql

import (
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserProfile struct {
	UserID       int
	OutputFormat int
	ValueFormat  int
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./telegram.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS user_profiles (
        user_id INTEGER PRIMARY KEY,
        outputformat INTEGER,
        valueformat INTEGER
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetUserProfile(db *sql.DB, userID int) (*UserProfile, error) {
	row := db.QueryRow("SELECT outputformat, valueformat FROM user_profiles WHERE user_id = ?", userID)

	var profile UserProfile
	err := row.Scan(&profile.OutputFormat, &profile.ValueFormat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}
	profile.UserID = userID
	return &profile, nil
}

func saveUserProfile(db *sql.DB, profile *UserProfile) error {
	_, err := db.Exec("INSERT OR REPLACE INTO user_profiles (user_id, outputformat, valueformat) VALUES (?, ?, ?)",
		profile.UserID, profile.OutputFormat, profile.ValueFormat)
	if err != nil {
		return err
	}

	return nil
}

func HandleCommand(db *sql.DB, update tgbotapi.Update) error {
	userID := update.Message.From.ID
	profile, err := GetUserProfile(db, userID)
	if err != nil {
		profile = &UserProfile{
			UserID:       userID,
			OutputFormat: 1,
			ValueFormat:  1,
		}
	}

	request := update.Message.Text

	switch request {
	case "/outputformat":
		profile.OutputFormat = 2
		saveUserProfile(db, profile)
	case "/valueformat":
		profile.ValueFormat = 2
		saveUserProfile(db, profile)
	case "/start":

	default:
		return fmt.Errorf("unknown command")
	}

	return nil
}
