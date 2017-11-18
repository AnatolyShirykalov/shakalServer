package hand

import (
	"github.com/googollee/go-socket.io"
	"log"
)

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
