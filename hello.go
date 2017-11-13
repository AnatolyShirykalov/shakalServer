package main

import (
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
	"reflect"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("action", func(f interface{}) {
			v := reflect.ValueOf(f)
			v.SetMapIndex(reflect.ValueOf("type"), reflect.ValueOf("aaaaa"))
			so.Emit("action", v.Interface())
			log.Println("on action", v)
		})
		so.On("disconnect", func() {
			log.Println("on disconnect")
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
