package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Events struct {
	Events []Event
}

func (events *Events) GetEventsForToday() *Events {
	eventsForToday := &Events{}

	for _, event := range events.Events {
		if event.Date.Format("01/02/2006") == time.Now().Format("01/02/2006") {
			eventsForToday.Events = append(eventsForToday.Events, event)
		}
	}

	return eventsForToday
}

func (events *Events) GetAllEventsOnFuture() *Events {
	eventsOnFuture := &Events{}

	for _, event := range events.Events {
		if time.Until(event.Date) > time.Hour*24 {
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

func GetEventsDb() (*Events, error) {
	events := &Events{}

	rows, err := db.Query("SELECT name, place, date FROM events")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		event := Event{}
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

func SetEventsDb(events *Events) error {
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

func main() {
	if err := InitDb(); err != nil {
		log.Fatal(err)
	}

	events, err := GetEventsDb()
	if err != nil {
		log.Fatal(err)
	}
	mode := ""

	reader := bufio.NewReader(os.Stdin)

	for mode != "exit" {
		fmt.Print("Input mode: ")
		fmt.Scanln(&mode)

		switch mode {
		case "add":
			var name, place, date string
			fmt.Print("Input name: ")
			fmt.Scanln(&name)
			fmt.Print("Input place: ")
			fmt.Scanln(&place)
			fmt.Print("Input date(01/02/2006 03:04): ")
			//fmt.Scanf("%s %s", &date, &time_s)
			date, err = reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			date = strings.TrimSpace(date)

			time, err := time.Parse("01/02/2006 03:04", date)
			if err != nil {
				log.Fatal(err)
			}
			events.Events = append(events.Events, Event{name, place, time})
		case "get":
			fmt.Println("All events: ")
			for _, event := range events.Events {
				fmt.Println(event)
			}

			fmt.Println("Events for today: ")
			for _, event := range events.GetEventsForToday().Events {
				fmt.Println(event)
			}

			fmt.Println("Events on future: ")
			for _, event := range events.GetAllEventsOnFuture().Events {
				fmt.Println(event)
			}

		default:
			fmt.Println("Unknown command.")
		}

		mode = ""
	}
}
