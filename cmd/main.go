package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type Events struct {
	Events []Event
}

func (events *Events) GetEventsForToday() *Events {
	eventsForToday := &Events{}

	for _, event := range events.Events {
		if event.Date.Format("02/01/2006") == time.Now().Format("02/01/2006") {
			eventsForToday.Events = append(eventsForToday.Events, event)
		}
	}

	return eventsForToday
}

func (events *Events) GetAllEventsOnFuture() *Events {
	eventsOnFuture := &Events{}

	for _, event := range events.Events {
		if event.Date.Sub(time.Now()) > time.Hour*24 {
			eventsOnFuture.Events = append(eventsOnFuture.Events, event)
		}
	}

	return eventsOnFuture
}

type Event struct {
	Name  string
	Place string
	Date  time.Time
}

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

func Close() error {
	return db.Close()
}

func main() {
	if err := InitDb(); err != nil {
		log.Fatal(err)
	}

	events := Events{
		Events: []Event{
			{
				Name:  "Go for a work",
				Place: "Office",
				Date:  time.Now(),
			},
			{
				Name:  "Go for a walk",
				Place: "Park",
				Date:  time.Now(),
			},
			{
				Name:  "Go for a work",
				Place: "Office",
				Date:  time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Name:  "Go for a walk",
				Place: "Park",
				Date:  time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Name:  "Go for a walk",
				Place: "Park",
				Date:  time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Name:  "Go for a library",
				Place: "Library",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Name:  "Go for a club",
				Place: "Club",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	fmt.Println("Hello, I am simple notifier.")
	fmt.Println("Events for today:")
	fmt.Println(events.GetEventsForToday().Events)
	fmt.Println("All events on future:")
	fmt.Println(events.GetAllEventsOnFuture())
	fmt.Println("All events:")
	fmt.Println(events.Events)
}
