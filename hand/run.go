package hand

import (
	"log"
)

func (hand *Hand) run() {
	for message := range hand.Ch {
		log.Println("received", message.Type)
		switch message.Type {
		case "connection":
			hand.Connect(message.So)
		case "join":
			gameId, playerId, err := getJoinParams(message.Data)
			if err == nil {
				hand.Join(gameId, playerId)
			}
			message.So.Emit("kaka", "lala")
		case "action":
			hand.Action()
		case "disconnection":
			hand.Disconnect()
		}
	}
}
