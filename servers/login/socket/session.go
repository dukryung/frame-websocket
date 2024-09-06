package socket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type SessionService struct {
	Upgrader *websocket.Upgrader
	Clients  map[string]*Client
	mu       sync.Mutex
}

type Client struct {
	SessionId string `json:"session-id"`
	Conn      *websocket.Conn
	Data      map[string]interface{}
}

func NewSessionService() *SessionService {
	return &SessionService{
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// This allows any origin, modify this for production to restrict origins
				return true
			},
		},
	}
}

func (ss *SessionService) CreateSession(c *gin.Context) {
	conn, err := ss.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	//TODO : 회원 정보 체크 로직

	sessionId := c.Query("session-id")
	if sessionId == "" {
		sessionId = "test"
	}

	client := &Client{
		SessionId: sessionId,
		Conn:      conn,
		Data:      make(map[string]interface{}),
	}

	ss.AddClient(client)

	go handleMessage(client)

}

func (ss *SessionService) AddClient(client *Client) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.Clients[client.SessionId] = client
	log.Printf("Client %s connected", client.SessionId)

}

// RemoveClient removes a client from the session manager
func (ss *SessionService) RemoveClient(clientID string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	delete(ss.Clients, clientID)
	fmt.Printf("Client %s disconnected\n", clientID)
}
func handleMessage(client *Client) {
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received message from %s: %s\n", client.SessionId, msg)

		// Echo the message back to the client (for testing)
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}

}
