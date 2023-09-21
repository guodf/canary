package netutil

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Closed = errors.New("WebSocket is closed")
)

type Event struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

type wsocketClient struct {
	Url     string
	Header  *http.Header
	handers map[string]EventHandlerFunc
	conn    *websocket.Conn
	isConn  bool
}

type wsocketServer struct {
	w              http.ResponseWriter
	r              *http.Request
	conn           *websocket.Conn
	responseHeader *http.Header
	handers        map[string]EventHandlerFunc
	isConn         bool
	mutex          sync.Mutex
}

type EventHandlerFunc func(event *Event)

type IWSocket interface {
	Watcher(eventType string, handler EventHandlerFunc)
	Run() error
	Close()
	Send(event Event) error
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWSocketServer(w http.ResponseWriter, r *http.Request, responseHeader *http.Header) IWSocket {
	return &wsocketServer{
		conn:           nil,
		w:              w,
		r:              r,
		responseHeader: responseHeader,
		handers:        map[string]EventHandlerFunc{},
		isConn:         false,
	}
}

func (ws *wsocketServer) Watcher(eventType string, handler EventHandlerFunc) {
	ws.handers[eventType] = handler
}

func (ws *wsocketServer) Run() error {
	c, err := upgrader.Upgrade(ws.w, ws.r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return err
	}
	ws.conn = c
	ws.isConn = true
	ws.listen()
	return nil
}

func (ws *wsocketServer) listen() {
	defer func() {
		log.Println("close conn")
		if ws.conn != nil {
			ws.conn.Close()
		}
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		var event Event
		err := ws.conn.ReadJSON(&event)
		if err != nil {
			log.Println("read:", err)
			return
		}
		go ws.handers[event.EventType](&event)
	}
}

func (ws *wsocketServer) Close() {
	ws.isConn = false
	ws.conn.Close()
}

func (ws *wsocketServer) Send(event Event) error {
	if ws.isConn {
		ws.mutex.Lock()
		defer ws.mutex.Unlock()
		return ws.conn.WriteJSON(event)
	}
	return Closed
}
