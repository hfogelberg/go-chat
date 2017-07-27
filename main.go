package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func main() {
	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// get room going
	go r.run()

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Printf("Listen and Serve: %s\n", err.Error())
		return
	}
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}
