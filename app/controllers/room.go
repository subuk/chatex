package controllers

import (
	"chatex/app/chat"
	"chatex/app/models"
	"code.google.com/p/go.net/websocket"
	"fmt"
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

func (self Room) Publish(roomId int64, text string, image_url string) revel.Result {
	self.Validation.Required(text)
	if self.Validation.HasErrors() {
		return self.RenderJson("text required")
	}
	var room = models.GetRoom(roomId)
	if room == nil {
		return self.NotFound("Room was not found")
	}
	var msg = models.Message{
		Text:     text,
		ImageUrl: image_url,
	}
	room.AddMessage(msg)
	return self.RenderJson(true)
}

func (self Room) Subscribe(ws *websocket.Conn, roomId int64) revel.Result {
	var room = models.GetRoom(roomId)
	if room == nil {
		return self.NotFound("Room was not found")
	}

	sub := chat.Join(fmt.Sprintf("Anonymous[%s]", ws.Request().RemoteAddr), roomId)

	defer chat.Cancel(sub)
	defer ws.Close()

	for {
		event := <-sub.Channel
		revel.INFO.Printf("%s: Received chat event %s", sub.String(), event.String())

		if event.Type == chat.EVENT_MSG {
			revel.INFO.Printf("Publishing message to room #%d\n", event.RoomId)
			if websocket.JSON.Send(ws, event) != nil {
				revel.WARN.Println("Client disconnected")
				break
			}
		} else if event.Type == chat.EVENT_PING {
			if websocket.JSON.Send(ws, event) != nil {
				revel.WARN.Println("Client disconnected")
				break
			}
		} else if event.Type == chat.EVENT_NEW_USER {
			websocket.JSON.Send(ws, event)
		} else if event.Type == chat.EVENT_USER_DISCONNECTED {
			websocket.JSON.Send(ws, event)
		} else {
			revel.ERROR.Fatal("Unknown event", event)
		}

	}

	return nil
}
