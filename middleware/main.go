package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"github.com/gorilla/websocket"
	"text/template"
)

const (
	srvAddr         = "224.0.0.1:9999"
	maxDatagramSize = 8192
	httpServiceAddr	= "localhost:8080"
)

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


func onStockDataReceived(src *net.UDPAddr, n int, message []byte) {
	data := []byte(message[:n])
	h.broadcast <- data
}

func serveMulticastUDP(a string, h func(*net.UDPAddr, int, []byte)) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}
}



func homeHandler(c http.ResponseWriter, req *http.Request) {
    homeTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../web"))))
	homeTempl = template.Must(template.ParseFiles("../web/client.html"))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)

	go h.run()

	//启动监听，接收股票数据
	go serveMulticastUDP(srvAddr, onStockDataReceived)

	//检查连接存活性
	go h.checkConnection()

	//监听客户端连接
	log.Printf("Listening to %s ...", httpServiceAddr)
	if err := http.ListenAndServe(httpServiceAddr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}





