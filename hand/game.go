package hand

import (
	"log"
)

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
