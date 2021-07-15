package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグ解釈
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run() //チャットルームを開始

	log.Println("Webサーバを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil { //Webサーバを起動
		log.Fatal("ListenAndServe:", err)
	}
}