package chat

import (
	"container/list"
	"github.com/robfig/revel"
	"time"
)

const (
	PING_TIMEOUT            = 3
	ALL_ROOMS               = 0
	EVENT_MSG               = 1
	EVENT_NEW_USER          = 2
	EVENT_PING              = 3
	EVENT_USER_DISCONNECTED = 4
)

var EVENT_NAMES = map[int]string{
	EVENT_PING:              "Ping",
	EVENT_MSG:               "NewMessage",
	EVENT_NEW_USER:          "NewUserConnected",
	EVENT_USER_DISCONNECTED: "UserDisconnected",
}

type Event struct {
	Type    int
	RoomId  int64
	Payload interface{}
}

func (self *Event) String() string {
	return EVENT_NAMES[self.Type]
}

type Subscription struct {
	Channel chan *Event
	Name    string
	RoomId  int64
	el      *list.Element
}

func (self *Subscription) String() string {
	return self.Name
}

var Subscribers = list.New()

func Publish(Type int, RoomId int64, Payload interface{}) {
	evt := Event{
		Type:    Type,
		RoomId:  RoomId,
		Payload: Payload,
	}

	for el := Subscribers.Front(); el != nil; el = el.Next() {
		s := el.Value.(Subscription)
		if evt.RoomId == s.RoomId || evt.RoomId == ALL_ROOMS {
			revel.INFO.Printf("--> Publishing event %s to subscription %s", evt.String(), s.String())
			s.Channel <- &evt
			revel.INFO.Printf("--> Event %s published to subscription %s", evt.String(), s.String())
		}
	}
}

func Join(name string, roomId int64) Subscription {
	var s = Subscription{
		Channel: make(chan *Event, 5),
		Name:    name,
		RoomId:  roomId,
	}

	s.el = Subscribers.PushFront(s)
	revel.INFO.Printf("Joining chat. Now: %d", Subscribers.Len())
	payload := map[string]interface{}{
		"Username": s.Name,
		"Total":    Subscribers.Len() + 1,
	}
	Publish(EVENT_NEW_USER, roomId, payload)
	return s
}

func Cancel(s Subscription) {
	revel.INFO.Printf("Canceling chat subscription. Now: %d", Subscribers.Len()-1)
	payload := map[string]interface{}{
		"Username": s.Name,
		"Total":    Subscribers.Len() - 1,
	}
	Publish(EVENT_USER_DISCONNECTED, s.RoomId, payload)
	Subscribers.Remove(s.el)
}

func pingSender() {
	for {
		time.Sleep(PING_TIMEOUT * time.Second)
		Publish(EVENT_PING, ALL_ROOMS, nil)
	}
}

func init() {
	go listenMessages()
	go pingSender()
}
