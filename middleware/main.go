package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"text/template"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
} 

var homeTempl *template.Template

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	clientAddr := r.RemoteAddr
	log.Printf("new connection. addr: %s", clientAddr)


	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writer()
	go c.reader()
}

func broadcast() {
	timer := time.NewTicker(100 * time.Millisecond)

	for {
		for {
			message, empty := popMessage(); if empty {
				break
			}

			h.broadcast <- []byte(message)
			
		}
		
		<-timer.C
	}
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
    homeTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../web"))))
	homeTempl = template.Must(template.ParseFiles("../web/client.html"))
	go h.run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	log.Printf("Listening to %s ...", *addr)
	go broadcast()
	go h.checkConnection()
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}