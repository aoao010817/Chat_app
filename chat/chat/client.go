package main
import (
	"github.com/gorilla/websocket"
	"time"
)
//チャットを行っている一人のユーザーを表す
type client struct {
	socket *websocket.Conn //このクライアントのためのWebsocket
	send chan *message //メッセージが送られるチャネル
	room *room //このクライアントが参加しているチャットルーム
	userData map[string]interface{} // ユーザーデータを保持
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			// if AvatarURL, ok := c.userData["avatar_url"]; ok {
			// 	msg.AvatarURL = AvatarURL.(string)
			// }
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}