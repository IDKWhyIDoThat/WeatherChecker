package dbsql

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbfilename = "./telegram.db"

type UserProfile struct {
	UserID       int
	OutputFormat int
	ValueFormat  int
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbfilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
			profile := &UserProfile{
				UserID:       userID,
				OutputFormat: 1,
				ValueFormat:  1,
			}
			log.Print("NewProfile")
			return profile, nil
		}
		return nil, err
	}
	profile.UserID = userID
	return &profile, nil
}

func SaveUserProfile(db *sql.DB, profile *UserProfile) error {
	_, err := db.Exec("INSERT OR REPLACE INTO user_profiles (user_id, outputformat, valueformat) VALUES (?, ?, ?)",
		profile.UserID, profile.OutputFormat, profile.ValueFormat)
	if err != nil {
		return err
	}
	return nil
}
