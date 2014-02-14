package chat

import (
	"chatex/app/models"
)

const (
	EVENT_MSG  = 1
	EVENT_USER = 2
)

type Event struct {
	Type    int
	Payload interface{}
}

type Group struct {
	Room     *models.Room
	Messages chan *Event
}

func MakeGroup(room *models.Room) Group {
	group := Group{
		Room: room,
	}
	group.Messages = make(chan *Event, 10)
	return group
}

func (self *Group) Join(id int64) {

}

func (self *Group) Exit(id int64) {

}

var EventChannel = make(chan *Event, 10)

func init() {
	go listenMessages()
}
