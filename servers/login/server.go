package login

import (
	"github.com/dukryung/frame-websocket/servers/login/socket"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	route          *gin.Engine
	sessionService *socket.SessionService
}

func NewServer() *Server {
	router := gin.Default()
	return &Server{
		route:          router,
		sessionService: socket.NewSessionService(),
	}
}

func (s *Server) Run() error {

	s.route.GET("/ws", s.sessionService.CreateSession)

	err := http.ListenAndServe(":8080", s.route)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() {

}
