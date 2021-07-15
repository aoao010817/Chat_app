package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"golang/trace"
)


type room struct {
	forward chan []byte      //他のクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     //参加しようとしているクライアント
	leave   chan *client     //退室しようとしているクライアント
	clients map[*client]bool //在室しているすべてのクライアント
}

func newRoom() *room { //すぐに利用できるチャットルームを返す
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer trace.Tracer
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true //参加
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send) //退室
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました")
			for client := range r.clients { //すべてのクライアントにメッセージを送信
				select {
				case client.send <- msg: //メッセージを送信
					r.tracer.Trace("クライアントに送信されました")
				default:
					delete(r.clients, client) //送信に失敗
					close(client.send)
					r.tracer.Trace("メッセージの送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("SeaveHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() {r.leave <- client}()
	go client.write()
	client.read()
}