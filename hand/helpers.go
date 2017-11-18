package hand

import (
	"errors"
	"log"
)

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
