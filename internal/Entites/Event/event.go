package event

import (
	"time"
)

type Events struct {
	Events []Event
}

type Event struct {
	Name  string
	Place string
	Date  time.Time
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
