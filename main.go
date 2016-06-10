package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"path/filepath"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}
func getMyEnv() map[string]string {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return env
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	env := getMyEnv() //ローカル環境変数読み込み
	flag.Parse()      // フラグを解釈します。
	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey(env["SEC_KEY"])
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の鍵", "http://localhost:8080/auth/callback/github"),
		google.New(env["GOOGLE_CLIENT_ID"], env["GOOGLE_ACCESS_KEY"], "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom(UseGravatar)
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	//チャットルームを開始します。
	go r.run()
	//ウェブサーバーを開始します
	log.Println("webサーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
