package controllers

import (
	"chatex/app/chat"
	"chatex/app/models"
	"chatex/app/routes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/robfig/revel"
)

type Room struct {
	*revel.Controller
}

func (self Room) History(roomId int64) revel.Result {
	var room = models.GetRoom(roomId)
	if room == nil {
		return self.NotFound("Room was not found")
	}
	return self.RenderJson(room.GetMessages())
}

func (self Room) Publish(roomId int64, text string) revel.Result {
	var room = models.GetRoom(roomId)
	if room == nil {
		return self.NotFound("Room was not found")
	}
	var msg = models.Message{
		Text: text,
	}
	room.AddMessage(msg)
	return self.Redirect(routes.Room.History(roomId))
}

func (self Room) Subscribe(ws *websocket.Conn, roomId int64) revel.Result {
	var room = models.GetRoom(roomId)
	var chatGroup = chat.MakeGroup(room)
	// for {
	// 	var msg models.Message
	// 	err := websocket.JSON.Receive(ws, &msg)
	// 	msg.RoomId = roomId
	// 	if err != nil {
	// 		fmt.Printf("Websocket message parse error: %s\n", err)
	// 		return nil
	// 	}
	// 	if msg.Text == "" {
	// 		fmt.Println("Blang message body")
	// 		continue
	// 	}
	// 	fmt.Printf("Publishing message with text %s to room #%d\n", msg.Text, msg.RoomId)
	// 	chatGroup.Messages <- &chat.Event{
	// 		Type:    chat.EVENT_MSG,
	// 		Payload: msg,
	// 	}
	// }
	for {
		select {
		case msg := <-chatGroup.Messages:
			fmt.Printf("Publishing message with text %s to room #%d\n", msg.Text, msg.RoomId)
		}
	}
	ws.Close()
	return nil
}
