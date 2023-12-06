package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/signature"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	addr := flag.String("addr", ":8080", "server port")

	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(github.New(clientId,
		clientSecret,
		"http://localhost:8080/auth/callback/github"))

	r := newRoom()

	http.Handle("/chat", MustAuth(&templateHandler{
		filename: "chat.html",
	}))

	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/", loginHandler)

	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}
