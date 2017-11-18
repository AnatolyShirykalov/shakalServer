package hand

import (
	"github.com/googollee/go-socket.io"
)

type Message struct {
	Type string
	Data interface{}
	So   socketio.Socket
}
