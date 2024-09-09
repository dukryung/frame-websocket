package login

import (
	"context"
	"fmt"
	"github.com/dukryung/frame-websocket/servers/login/socket"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router         *gin.Engine
	ctx            context.Context
	sessionService *socket.SessionService
	srv            *http.Server
	close          chan bool
}

func NewServer() *Server {
	router := gin.Default()
	return &Server{
		router:         router,
		ctx:            context.Background(),
		sessionService: socket.NewSessionService(),
		srv: &http.Server{
			Addr:    "127.0.0.1:12345",
			Handler: router,
		},
		close: make(chan bool),
	}
}

func (s *Server) Run() error {

	s.router.GET("/ws", s.sessionService.CreateSession)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v\n", err)
		}
	}()

	<-s.close

	err := s.srv.Shutdown(s.ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Server) Close() {
	s.close <- true
}
