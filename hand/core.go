package hand

import (
	"github.com/googollee/go-socket.io"
	"log"
	"reflect"
	"strings"
)

func (h *Hand) HasGame(id string) bool {
	_, ok := h.Games[id]
	return ok
}

func (h *Hand) GetGame(id string) *Game {
	g, _ := h.Games[id]
	return g
}

func (h *Hand) CreateGame(id string) (g *Game) {
	log.Println("hand create game", id)
	g = &Game{}
	h.Games[id] = g
	return
}

func (h *Hand) NextCounter() {
	h.counter += 1
	if h.counter == 10000000 {
		h.counter = 1
	}
}

func (h *Hand) Connect(so socketio.Socket) {
	log.Println("hand connect")
	h.NextCounter()
	so.Emit("free id", h.counter)
}

func (h *Hand) Join(gameId, playerId string) map[string]string {
	log.Println("hand join")
	if !h.HasGame(gameId) {
		game := h.CreateGame(gameId)
		index := game.AddPlayer(playerId)
		game.ConnectPlayer(playerId)
		return map[string]string{"playerId": playerId, "index": index}
	}
	game := h.GetGame(gameId)
	if len(game.Players) < 4 {
		index := game.AddPlayer(playerId)
		game.ConnectPlayer(playerId)
		return map[string]string{"playerId": playerId, "index": index}
	}
	if game.ConnectedPlayers() < 4 {
		index, player := game.GetPlayer(playerId)
		if player == nil {
			return map[string]string{"watcher": "true"}
		}
		game.ConnectPlayer(playerId)
		return map[string]string{"playerId": playerId, "index": index}
	}
	return map[string]string{"watcher": "true"}
}

func (h *Hand) Action(data interface{}, s socketio.Socket) {
	log.Println("hand Action")
	v := reflect.ValueOf(data)
	typeWithPrefix, ok := v.Interface().(map[string]interface{})["type"]
	if !ok {
		log.Println(v)
	}
	gameId, ok := v.Interface().(map[string]interface{})["gameId"]
	if !ok {
		log.Println(v)
	}
	tp := strings.Replace(typeWithPrefix.(string), "server/", "", 1)
	v.SetMapIndex(reflect.ValueOf("type"), reflect.ValueOf(tp))
	ret := v.Interface()
	s.Emit("action", ret)
	s.BroadcastTo(gameId.(string), "action", ret)
}

func (h *Hand) Disconnect() {
	log.Println("hand disconnect")
}

func (h *Hand) ChatMessage() {
	log.Println("hand chat message")
}
