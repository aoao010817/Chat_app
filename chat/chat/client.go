package main
import (
	"github.com/gorilla/websocket"
)
//チャットを行っている一人のユーザーを表す
type client struct {
	socket *websocket.Conn //このクライアントのためのWebsocket
	send chan []byte //メッセージが送られるチャネル
	room *room //このクライアントが参加しているチャットルーム
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}