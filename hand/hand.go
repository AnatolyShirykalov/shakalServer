package hand

import (
	"github.com/googollee/go-socket.io"
)

type Hand struct {
	Ch      chan Message
	Games   map[string]*Game
	counter int64
}

func NewHand() (ret *Hand) {
	ret = &Hand{}
	ret.Ch = make(chan Message)
	ret.Games = make(map[string]*Game)
	go ret.run()
	return ret
}

func (hand *Hand) OnAction(s socketio.Socket, data interface{}) {
	hand.Ch <- Message{Type: "action", Data: data, So: s}
	/*v := reflect.ValueOf(data)*/
	//typeWithPrefix, ok := v.Interface().(map[string]interface{})["type"]
	//if !ok {
	//log.Println(v)
	//}
	//tp := strings.Replace(typeWithPrefix.(string), "server/", "", 1)
	//v.SetMapIndex(reflect.ValueOf("type"), reflect.ValueOf(tp))
	//ret := v.Interface()
	//s.Emit("action", ret)
	/*s.BroadcastTo("chat", "action", ret)*/
}

func (hand *Hand) OnJoin(s socketio.Socket, data interface{}) {
	hand.Ch <- Message{Type: "join", Data: data, So: s}
	//s.Join(roomId)
}
