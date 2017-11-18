package hand

import (
	"log"
	"strconv"
)

type Game struct {
	Players []*Player
	Room    string
}

func (g *Game) AddPlayer(id string) string {
	log.Println("game add player", id)
	g.Players = append(g.Players, &Player{Id: id, Connected: true})
	return strconv.Itoa(len(g.Players) - 1)
}

func (g *Game) HasPlayer(id string) bool {
	log.Println("game has player", id)
	for _, player := range g.Players {
		if player.Id == id {
			return true
		}
	}
	return false
}

func (g *Game) ConnectPlayer(id string) {
	log.Println("game connect player", id)
	for _, player := range g.Players {
		if player.Id == id {
			player.Connected = true
		}
	}
}

func (g *Game) ConnectedPlayers() (ret int) {
	for _, player := range g.Players {
		if player.Connected {
			ret += 1
		}
	}
	return
}

func (g *Game) GetPlayer(id string) (string, *Player) {
	for index, player := range g.Players {
		if player.Id == id {
			return strconv.Itoa(index), player
		}
	}
	return "-1", nil
}
