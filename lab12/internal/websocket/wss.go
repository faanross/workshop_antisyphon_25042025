package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var WebSocketPort = 8080

type WebSocketServer struct {
	port     int
	client   *websocket.Conn
	clientMx sync.Mutex
}

func NewWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{
		port:   port,
		client: nil,
	}
}

var GlobalWSServer *WebSocketServer

func StartWebSocketServer() {
	GlobalWSServer = NewWebSocketServer(WebSocketPort)

	go func() {
		err := GlobalWSServer.Start()
		if err != nil {
			log.Fatalf("WebSocket server failed to start: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("WebSocket server started on port ", WebSocketPort)
}

func (s *WebSocketServer) Start() error {
	http.HandleFunc("/ws", s.handleWebSocket)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("WebSocket server listening on %s\n", addr)

	return http.ListenAndServe(addr, nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrader failed to upgrade: %v", err)
		return
	}
	defer func() {
		conn.Close()

		s.clientMx.Lock()
		s.client = nil
		s.clientMx.Unlock()
	}()

	s.clientMx.Lock()
	s.client = conn
	s.clientMx.Unlock()

	fmt.Println("WebSocket connection established")

	for {
		var msg Message

		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket connection failed to read: %v", err)
			break
		}

		log.Printf("WebSocket connection read: %+v", msg)

		if msg.Type == CommandMessage {
			AgentCommands.QueueCommand(msg.Command)
		}
	}

}

func (s *WebSocketServer) Broadcast(msg Message) {
	s.clientMx.Lock()
	defer s.clientMx.Unlock()

	if s.client != nil {
		err := s.client.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to broadcast: %v", err)
			return
		}
	}
}
