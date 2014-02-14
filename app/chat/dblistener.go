package chat

import (
	"chatex/app/models"
	"fmt"
	"github.com/lib/pq"
	"github.com/robfig/revel"
	"strconv"
	"strings"
	"time"
)

func listenerEventsHandler(event pq.ListenerEventType, err error) {
	if err != nil {
		panic(err)
	}
	if event == pq.ListenerEventConnected {
		fmt.Println("Database listener successfully connected")
	}

}

func listenMessages() {
	fmt.Println("Setting up db listening")
	var Spec, _ = revel.Config.String("db.spec")
	var listener = pq.NewListener(Spec, 10*time.Second, time.Minute, listenerEventsHandler)
	var err = listener.Listen("new_message")
	if err != nil {
		panic(err)
	}
	waitForNotification(listener)
}

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case notify := <-l.Notify:

			payload := strings.SplitN(notify.Extra, "|", 3)

			id, err := strconv.ParseInt(payload[0], 10, 64)
			if err != nil {
				panic(err)
			}
			roomId, err := strconv.ParseInt(payload[1], 10, 64)
			if err != nil {
				panic(err)
			}
			msg := models.Message{
				Id:     id,
				RoomId: roomId,
				Text:   payload[2],
			}
			fmt.Printf("received notification with payload: '%d' '%d' '%s'\n", msg.Id, msg.RoomId, msg.Text)
			EventChannel <- &Event{
				Type:    EVENT_MSG,
				Payload: msg,
			}

		case <-time.After(200 * time.Millisecond):
			go func() {
				if err := l.Ping(); err != nil {
					panic(err)
				}
			}()
		}
	}
}
