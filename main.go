package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"src/github.com/stretchr/gomniauth"
	"src/github.com/stretchr/gomniauth/providers/facebook"
	"src/github.com/stretchr/gomniauth/providers/github"
	"src/github.com/stretchr/gomniauth/providers/google"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈します。
	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("zMsQ7B7FQneKKXTQKviECsjHjvA%RNKA")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/github"),
		google.New("1028664644457-5gr0979l96j6ku8a2lfrm3id2ph5dvo0.apps.googleusercontent.com",
			"JXCqNsv0aWeh6p9jv81lCpcx", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//チャットルームを開始します。
	go r.run()
	//ウェブサーバーを開始します
	log.Println("webサーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
