package hand

import (
	"github.com/googollee/go-socket.io"
	"log"
)

func Server() *socketio.Server {
	h := NewHand()
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		so.On("free id", func(data interface{}) {
			h.Ch <- Message{Type: "connect", So: so}
		})
		so.Join("chat")
		so.On("chat message", func(msg string) {
			h.Ch <- Message{Type: "chat message", So: so}
			//log.Println("emit:", so.Emit("chat message", msg))
			//so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("action", h.OnAction)
		so.On("join", h.OnJoin)
		so.On("disconnection", func() {
			h.Ch <- Message{Type: "disconnect", So: so}
			log.Println("on disconnect")
			//for _, roomId := range so.Rooms() {
			//so.BroadcastTo(roomId, "leaving", so.Id)
			//}
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	return server
}
