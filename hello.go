package main

import (
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
	"reflect"
	"strings"
)

type Message struct {
	Type string
	Data interface{}
	so   socketio.Socket
}

type Player struct {
	Id        string
	Connected bool
}

type Game struct {
	Players []Player
	Room    string
}

type Hand struct {
	Ch    chan Message
	Games map[string]Game
}

func NewHand() (ret *Hand) {
	ret = &Hand{}
	ret.Ch = make(chan Message)
	go func(ch chan Message) {
		for message := range ch {
			log.Println("received", message.Type)
		}
	}(ret.Ch)
	return ret
}

func (hand *Hand) onAction(s socketio.Socket, data interface{}) {
	hand.Ch <- Message{Type: "action", Data: data, so: s}
	v := reflect.ValueOf(data)
	typeWithPrefix, ok := v.Interface().(map[string]interface{})["type"]
	if !ok {
		log.Println(v)
	}
	tp := strings.Replace(typeWithPrefix.(string), "server/", "", 1)
	v.SetMapIndex(reflect.ValueOf("type"), reflect.ValueOf(tp))
	ret := v.Interface()
	s.Emit("action", ret)
	s.BroadcastTo("chat", "action", ret)
}

func (hand *Hand) onJoin(s socketio.Socket, roomId string) {
	hand.Ch <- Message{Type: "join", Data: roomId, so: s}
	s.Join(roomId)
}

func main() {
	hand := NewHand()
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		hand.Ch <- Message{Type: "connection", so: so}
		so.Join("chat")
		so.On("chat message", func(msg string) {
			hand.Ch <- Message{Type: "chat message", so: so}
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("action", hand.onAction)
		so.On("join", hand.onJoin)
		so.On("disconnect", func() {
			hand.Ch <- Message{Type: "disconnect", so: so}
			log.Println("on disconnect")
			for _, roomId := range so.Rooms() {
				so.BroadcastTo(roomId, "leaving", so.Id)
			}
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
