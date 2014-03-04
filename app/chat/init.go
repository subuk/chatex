package chat

import (
	"time"
)

const (
	PING_TIMEOUT = 3
	ALL_ROOMS    = 0
)

const (
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
