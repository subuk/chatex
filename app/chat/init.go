package chat

import (
	"container/list"
	"github.com/robfig/revel"
	"time"
)

const (
	EVENT_MSG      = 1
	EVENT_NEW_USER = 2
)

type Event struct {
	Type    int
	Payload interface{}
}

var Subscribers = list.New()

func Publish(evt *Event) {
	for el := Subscribers.Front(); el != nil; el = el.Next() {
		s := el.Value.(chan *Event)
		s <- evt
	}
}

func Join() (chan *Event, *list.Element) {
	var s = make(chan *Event)
	el := Subscribers.PushFront(s)
	revel.INFO.Printf("Joining chat. Now: %d", Subscribers.Len())
	return s, el
}

func Cancel(el *list.Element) {
	Subscribers.Remove(el)
	revel.INFO.Printf("Canceling chat subscription. Now: %d", Subscribers.Len())
}

func subCountLogger() {
	for {
		time.Sleep(1 * time.Second)
		revel.TRACE.Printf("Subscribers count: %d\n", Subscribers.Len())
	}
}

func init() {
	go listenMessages()
	// go subCountLogger()
}
