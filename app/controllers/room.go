package controllers

import (
	"chatex/app/chat"
	"chatex/app/models"
	"code.google.com/p/go.net/websocket"
	"github.com/robfig/revel"
)

type Room struct {
	*revel.Controller
}

func (self *Room) Show(roomId int64) revel.Result {
	var room = models.GetRoom(roomId)
	if room == nil {
		return self.NotFound("Room was not found")
	}
	return self.Render(room)
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
	return self.RenderJson(true)
}

func (self Room) Subscribe(ws *websocket.Conn, roomId int64) revel.Result {
	var room = models.GetRoom(roomId)
	if room == nil {
		return self.NotFound("Room was not found")
	}

	eventChannel, el := chat.Join()
	defer chat.Cancel(el)

	for {
		e := <-eventChannel

		if e.Type == chat.EVENT_MSG {
			var msg = e.Payload.(models.Message)
			if msg.RoomId == room.Id {
				revel.INFO.Printf("Publishing message with text %s to room #%d\n", msg.Text, msg.RoomId)
				if websocket.JSON.Send(ws, msg) != nil {
					revel.WARN.Println("Client disconnected")
					break
				}
			}
		} else if e.Type == chat.EVENT_NEW_USER {
			if e.Payload.(int64) == room.Id {
				revel.INFO.Printf("New anoymous joined room #%d\n", room.Id)
				msg := models.Message{
					Text: "New user joined!",
				}
				if websocket.JSON.Send(ws, msg) != nil {
					revel.WARN.Println("Client disconnected")
					break
				}
			}
		}

	}

	ws.Close()
	return nil
}
