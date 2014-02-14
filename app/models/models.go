package models

import (
	"github.com/robfig/revel/modules/db/app"
	"time"
)

type Message struct {
	Id         int64
	RoomId     int64
	CreateTime time.Time
	Text       string
}

type Room struct {
	Id int64
}

func GetRoom(id int64) *Room {
	var count int64 = 0
	var room *Room = nil
	db.Db.QueryRow("SELECT COUNT(id) FROM room WHERE id=$1", id).Scan(&count)
	if count != 0 {
		room = &Room{Id: id}
	}
	return room
}

func (self *Room) GetMessages() []Message {
	var messages []Message
	var msg Message

	rows, err := db.Db.Query("SELECT id, text, create_time FROM message WHERE room_id = $1", self.Id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		msg = Message{RoomId: self.Id}
		rows.Scan(&msg.Id, &msg.Text, &msg.CreateTime)
		messages = append(messages, msg)
	}

	return messages
}

func (self *Room) AddMessage(message Message) {
	message.RoomId = self.Id
	row, err := db.Db.Query("INSERT INTO message (room_id, text) VALUES ($1, $2) RETURNING id, create_time", message.RoomId, message.Text)
	if err != nil {
		panic(err)
	}
	row.Scan(&message.Id, &message.CreateTime)
}
