package main

import (
	"./hand"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
	//"reflect"
	//"strings"
)

func main() {
	h := hand.NewHand()
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		h.Ch <- hand.Message{Type: "connection", So: so}
		so.Join("chat")
		so.On("chat message", func(msg string) {
			h.Ch <- hand.Message{Type: "chat message", So: so}
			//log.Println("emit:", so.Emit("chat message", msg))
			//so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("action", h.OnAction)
		so.On("xjoin", h.OnJoin)
		so.On("disconnection", func() {
			h.Ch <- hand.Message{Type: "disconnect", So: so}
			log.Println("on disconnect")
			//for _, roomId := range so.Rooms() {
			//so.BroadcastTo(roomId, "leaving", so.Id)
			//}
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.Static("/js", ".")

	r.GET("/socket.io/", gin.WrapH(server))

	r.Run("localhost:12312")
}
