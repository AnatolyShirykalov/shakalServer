package hand

import (
	"github.com/googollee/go-socket.io"
)

type Hand struct {
	Ch             chan Message
	Games          map[string]*Game
	counter        int64
	socketIdGameId map[string]string
}

func NewHand() (ret *Hand) {
	ret = &Hand{}
	ret.Ch = make(chan Message)
	ret.Games = make(map[string]*Game)
	ret.socketIdGameId = make(map[string]string)
	go ret.run()
	return ret
}

func (h *Hand) OnAction(s socketio.Socket, data interface{}) {
	h.Ch <- Message{Type: "action", Data: data, So: s}

}

func (h *Hand) OnJoin(s socketio.Socket, data interface{}) {
	h.Ch <- Message{Type: "join", Data: data, So: s}
	//s.Join(roomId)
}
