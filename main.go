package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth");err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグ解釈
	//Gomniauth
	gomniauth.SetSecurityKey("aoao010817")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の鍵", "http:localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の鍵", "http:localhost:8080/auth/callback/github"),
		google.New("62410808419-ttid6l47fc1gqvfb11flpjui9jrgd6l0.apps.googleusercontent.com", "TguyYcJDEYZXaHJfFNglSy9E", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/",http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars/"))))
	go r.run() //チャットルームを開始

	log.Println("Webサーバを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil { //Webサーバを起動
		log.Fatal("ListenAndServe:", err)
	}
}

// clientID: 62410808419-ttid6l47fc1gqvfb11flpjui9jrgd6l0.apps.googleusercontent.com
// seclet: TguyYcJDEYZXaHJfFNglSy9E