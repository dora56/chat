package main

import (
	"flag"
	"github.com/dora56/chat/trace"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHendler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHendler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHendler{filename: "chat.html"})
	http.Handle("/room", r)
	//チャットルームを開始します。
	go r.run()
	//ウェブサーバーを開始します
	log.Println("webサーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
