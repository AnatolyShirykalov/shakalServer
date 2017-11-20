package hand

import (
	"github.com/googollee/go-socket.io"
	"log"
	"reflect"
	"strings"
)

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

func (h *Hand) Join(gameId, playerId, socketId string) map[string]string {
	log.Println("hand join")
	h.socketIdGameId[socketId] = gameId
	if !h.HasGame(gameId) {
		game := h.CreateGame(gameId)
		index := game.AddPlayer(playerId)
		game.ConnectPlayer(playerId, socketId)
		return map[string]string{"playerId": playerId, "index": index}
	}
	game := h.GetGame(gameId)
	if len(game.Players) < 4 {
		index := game.AddPlayer(playerId)
		game.ConnectPlayer(playerId, socketId)
		return map[string]string{"playerId": playerId, "index": index}
	}
	if game.ConnectedPlayers() < 4 {
		index, player := game.GetPlayer(playerId)
		if player == nil {
			return map[string]string{"watcher": "true"}
		}
		game.ConnectPlayer(playerId, socketId)
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
	gameId, ok := h.socketIdGameId[s.Id()]
	if !ok {
		log.Println("cannot find game for socketId", s.Id)
		return
	}
	tp := strings.Replace(typeWithPrefix.(string), "server/", "", 1)
	v.SetMapIndex(reflect.ValueOf("type"), reflect.ValueOf(tp))
	ret := v.Interface()
	s.Emit("action", ret)
	s.BroadcastTo(gameId, "action", ret)
}

func (h *Hand) Disconnect(socketId string) (gameId, index string) {
	gameId, ok := h.socketIdGameId[socketId]
	if !ok {
		log.Println("cannot find game for socketId", socketId)
		return
	}
	delete(h.socketIdGameId, gameId)
	index, _ = h.Games[gameId].DisconnectPlayer(socketId)
	log.Println("hand disconnect")
	return
}

func (h *Hand) ChatMessage() {
	log.Println("hand chat message")
}
