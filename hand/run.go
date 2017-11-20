package hand

import (
	"log"
)

func (h *Hand) run() {
	for message := range h.Ch {
		log.Println("received", message.Type, message.So.Id)
		switch message.Type {
		case "connect":
			h.Connect(message.So)
		case "action":
			actionType := actionType(message.Data)
			log.Println(actionType)
			if actionType != "server/JOIN" {
				h.Action(message.Data, message.So)
				break
			}
			gameId, playerId, err := getJoinParams(message.Data)
			if err != nil {
				message.So.Emit("error", err.Error())
				break
			}
			msg := h.Join(gameId, playerId, message.So.Id())
			message.So.Join(gameId)
			msg["type"] = "JOIN"
			message.So.BroadcastTo(gameId, "action", msg)
			msg["you"] = "true"
			message.So.Emit("action", msg)
		case "disconnect":
			log.Println("message disconnection")
			gameId, index := h.Disconnect(message.So.Id())
			msg := map[string]string{"type": "disconnect", "index": index}
			message.So.BroadcastTo(gameId, "action", msg)
		}
	}
}
