package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	event "github.com/kasariks/simple_notifier/internal/Entites/Event"
	"github.com/kasariks/simple_notifier/internal/db"
	_ "modernc.org/sqlite"
)

func main() {
	if err := db.InitDb(); err != nil {
		log.Fatal(err)
	}

	events, err := db.GetEventsDb()
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
			date, err = reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			date = strings.TrimSpace(date)

			time, err := time.Parse("01/02/2006 03:04", date)
			if err != nil {
				log.Fatal(err)
			}
			events.Events = append(events.Events, event.Event{
				Name:  name,
				Place: place,
				Date:  time})
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
