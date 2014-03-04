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
	ImageUrl   string
}

type Room struct {
	Id          int64
	BoardId     int64
	Description string
}

func GetRoom(id int64) *Room {
	row := db.Db.QueryRow("SELECT board_id, description FROM room WHERE id=$1", id)
	if row == nil {
		return nil
	}
	room := Room{Id: id}
	row.Scan(&room.BoardId, &room.Description)

	return &room
}

func GetRoomList() []Room {
	var rooms []Room

	rows, err := db.Db.Query("SELECT id, board_id, description FROM room ORDER BY id")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var room = Room{}
		rows.Scan(&room.Id, &room.BoardId, &room.Description)
		rooms = append(rooms, room)
	}
	return rooms
}

func GetMessage(id int64) *Message {
	var msg Message
	row := db.Db.QueryRow("SELECT id, text, create_time, image_url FROM message WHERE id=$1", id)
	if row == nil {
		return nil
	}
	row.Scan(&msg.Id, &msg.Text, &msg.CreateTime, &msg.ImageUrl)
	return &msg
}

func (self *Room) GetMessages() []Message {
	var messages []Message
	rows, err := db.Db.Query("SELECT id, text, create_time, image_url FROM message WHERE room_id = $1 ORDER BY id", self.Id)
	if err != nil {
		panic(err)
	}

	var msg Message
	for rows.Next() {
		msg = Message{RoomId: self.Id}
		rows.Scan(&msg.Id, &msg.Text, &msg.CreateTime, &msg.ImageUrl)
		messages = append(messages, msg)
	}

	return messages
}

func (self *Room) AddMessage(message Message) {
	message.RoomId = self.Id
	row, err := db.Db.Query("INSERT INTO message (room_id, text, image_url) VALUES ($1, $2, $3) RETURNING id, create_time", message.RoomId, message.Text, message.ImageUrl)
	if err != nil {
		panic(err)
	}
	row.Scan(&message.Id, &message.CreateTime)
}
