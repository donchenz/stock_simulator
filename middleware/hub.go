package main

import  (
	"log"
	"time"
)

const checkConnectionInterval = 5
const checkConnectionExpire = 30

type hub struct {
	connections map[*connection]bool

	heartbeats map[*connection]int64

	broadcast chan []byte

	register chan *connection

	unregister chan *connection

	heartbeatRefresh chan * connection
}

var h = hub{
	broadcast:   		make(chan []byte),
	register:    		make(chan *connection),
	unregister:  		make(chan *connection),
	heartbeatRefresh :	make(chan *connection),
	connections: 		make(map[*connection]bool),
	heartbeats:			make(map[*connection]int64),
}

/*
检查连接存活性
*/
func (h *hub) checkConnection() {
	timer := time.NewTicker(5 * time.Second)

	for {
		for c := range h.heartbeats {
			tsBefore := h.heartbeats[c]
			tsNow := time.Now().Unix()
			diff := tsNow - tsBefore
			if diff > checkConnectionExpire {
				delete(h.heartbeats, c)
				h.unregister <- c
			}

		}
		
		<-timer.C
	}
}

func (h *hub) run() {
	for {
		select {
		/* 注册连接 */
		case c := <-h.register:
			log.Printf("register connection")
			h.connections[c] = true
		 /* 注销连接 */
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		/* 刷新心跳时间 */
		case c := <-h.heartbeatRefresh:
			h.heartbeats[c] = time.Now().Unix()
		/* 广播 */
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}

