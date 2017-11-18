package hand

import (
	"log"
)

func (h *Hand) run() {
	for message := range h.Ch {
		log.Println("received", message.Type)
		switch message.Type {
		case "connect":
			h.Connect(message.So)
		case "join":
			gameId, playerId, err := getJoinParams(message.Data)
			if err != nil {
				message.So.Emit("error", err.Error())
				break
			}
			msg := h.Join(gameId, playerId)
			message.So.Join(gameId)
			message.So.BroadcastTo(gameId, "joined", msg)
			message.So.Emit("joined", msg)
		case "action":
			h.Action(message.Data, message.So)
		case "disconnection":
			h.Disconnect()
		}
	}
}
