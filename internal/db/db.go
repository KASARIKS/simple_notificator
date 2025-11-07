package db

import (
	"database/sql"
	"os"
	"time"

	event "github.com/kasariks/simple_notifier/internal/Entites/Event"
)

const schema = "CREATE TABLE events (" +
	"id INTEGER PRIMARY KEY AUTOINCREMENT," +
	"name TEXT NOT NULL," +
	"place TEXT NOT NULL," +
	"date TEXT NOT NULL" +
	");"

var db *sql.DB

func InitDb() error {
	dbFile := "events.db"

	_, err := os.Stat(dbFile)

	var install = false
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}

func GetEventsDb() (*event.Events, error) {
	events := &event.Events{}

	rows, err := db.Query("SELECT name, place, date FROM events")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		event := event.Event{}
		strDate := ""
		rows.Scan(&event.Name, &event.Place, &strDate)
		event.Date, err = time.Parse("01/02/2006 03:04", strDate)
		if err != nil {
			return nil, err
		}

		events.Events = append(events.Events, event)
	}

	return events, nil
}

func SetEventsDb(events *event.Events) error {
	for _, event := range events.Events {
		_, err := db.Exec("INSERT INTO events (name, place, date) VALUES (:name, :place, :date)",
			sql.Named("name", event.Name),
			sql.Named("place", event.Place),
			sql.Named("date", event.Date.Format("01/02/2006 03:04")))
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	return db.Close()
}
