package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
	//"reflect"
	//"strings"
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
	Players [4]Player
	Room    string
}

func (g *Game) AddPlayer(id string) {
	log.Println("game add player", id)
}

func (g *Game) HasPlayer(id string) {
	log.Println("game has player", id)
}

func (g *Game) ConnectPlayer(id string) {
	log.Println("game connect player", id)
}

type Hand struct {
	Ch      chan Message
	Games   map[string]*Game
	counter int64
}

func getJoinParams(d interface{}) (string, string, error) {
	data := d.(map[string]interface{})
	gameId, ok := data["gameId"]
	if !ok {
		log.Fatal("has no gameId in message")
		return "", "", errors.New("getJoinParams::has no gameId in message")
	}
	playerId, ok := data["playerId"]
	if !ok {
		log.Fatal("has no playerId in message")
		return "", "", errors.New("getJoinParams::has no playerId in message")
	}
	return gameId.(string), playerId.(string), nil
}

func (hand *Hand) run() {
	for message := range hand.Ch {
		log.Println("received", message.Type)
		switch message.Type {
		case "connection":
			hand.Connect(message.so)
		case "join":
			// Data: {gameId, PlayerId}
			// gameId playerId
			gameId, playerId, err := getJoinParams(message.Data)
			if err == nil {
				hand.Join(gameId, playerId)
			}
			message.so.Emit("kaka", "lala")
		case "action":
			hand.Action()
		case "disconnection":
			hand.Disconnect()
		}
	}
}

func (hand *Hand) HasGame(id string) bool {
	_, ok := hand.Games[id]
	return ok
}

func (hand *Hand) CreateGame(id string) (g *Game) {
	log.Println("hand create game", id)
	g = &Game{}
	hand.Games[id] = g
	return
}

func (hand *Hand) NextCounter() {
	hand.counter += 1
	if hand.counter == 10000000 {
		hand.counter = 1
	}
}

func (hand *Hand) Connect(so socketio.Socket) {
	log.Println("hand connect")
	hand.NextCounter()
	so.Emit("free id", hand.counter)
}

func (hand *Hand) Join(gameId, playerId string) {
	log.Println("hand join")
	if !hand.HasGame(gameId) {
		game := hand.CreateGame(gameId)
		game.AddPlayer(playerId)
		game.ConnectPlayer(playerId)
		return
	}
}

func (hand *Hand) Action() {
	log.Println("hand Action")
}

func (hand *Hand) Disconnect() {
	log.Println("hand disconnect")
}

func (hand *Hand) ChatMessage() {
	log.Println("hand chat message")
}

func NewHand() (ret *Hand) {
	ret = &Hand{}
	ret.Ch = make(chan Message)
	ret.Games = make(map[string]*Game)
	go ret.run()
	return ret
}

func (hand *Hand) onAction(s socketio.Socket, data interface{}) {
	hand.Ch <- Message{Type: "action", Data: data, so: s}
	/*v := reflect.ValueOf(data)*/
	//typeWithPrefix, ok := v.Interface().(map[string]interface{})["type"]
	//if !ok {
	//log.Println(v)
	//}
	//tp := strings.Replace(typeWithPrefix.(string), "server/", "", 1)
	//v.SetMapIndex(reflect.ValueOf("type"), reflect.ValueOf(tp))
	//ret := v.Interface()
	//s.Emit("action", ret)
	/*s.BroadcastTo("chat", "action", ret)*/
}

func (hand *Hand) onJoin(s socketio.Socket, data interface{}) {
	hand.Ch <- Message{Type: "join", Data: data, so: s}
	//s.Join(roomId)
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
			//log.Println("emit:", so.Emit("chat message", msg))
			//so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("action", hand.onAction)
		so.On("xjoin", hand.onJoin)
		so.On("disconnection", func() {
			hand.Ch <- Message{Type: "disconnect", so: so}
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
