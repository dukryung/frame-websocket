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
	Sessions map[string]*Session
	mu       sync.Mutex
}

type Session struct {
	SessionId string `json:"session-id"`
	Conn      *websocket.Conn
	Data      map[string]interface{}
}

func NewSessionService() *SessionService {
	return &SessionService{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allows connections from all origins
			},
		},
		Sessions: make(map[string]*Session),
		mu:       sync.Mutex{},
	}
}

func (ss *SessionService) CreateSession(c *gin.Context) {
	conn, err := ss.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}

	sessionId := c.Query("session-id")
	if sessionId == "" {
		sessionId = "test"
	}

	session := &Session{
		SessionId: sessionId,
		Conn:      conn,
		Data:      make(map[string]interface{}),
	}

	ss.AddSession(session)

	go ss.HandleMessage(session)

}

func (ss *SessionService) AddSession(session *Session) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.Sessions[session.SessionId] = session
	fmt.Printf("Client %s connected\n", session.SessionId)

}

// RemoveSession removes a client from the session manager
func (ss *SessionService) RemoveSession(sessionID string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	delete(ss.Sessions, sessionID)
	fmt.Printf("Client %s disconnected\n", sessionID)
}
func (ss *SessionService) HandleMessage(session *Session) {
	defer ss.RemoveSession(session.SessionId)
	defer session.Conn.Close()

	for {
		_, msg, err := session.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Read Message error", err)
			return
		}

		fmt.Printf("Received message from %s: %s\n", session.SessionId, msg)
		// Echo the message back to the client (for testing)
		if err := session.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
