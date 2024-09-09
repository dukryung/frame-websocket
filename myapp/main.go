package main

import (
	"github.com/dukryung/frame-websocket/servers/login"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Run() error
	Close()
}

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverManager := NewServersManager(login.NewServer())

	serverManager.Run()
	<-quit
	serverManager.Close()

}

type ServersManger struct {
	servers []Server
}

func NewServersManager(servers ...Server) *ServersManger {
	return &ServersManger{
		servers: servers,
	}
}

func (sm *ServersManger) Run() {
	for _, server := range sm.servers {
		go func() {
			err := server.Run()
			if err != nil {
				panic(err)
			}
		}()
	}
}

func (sm *ServersManger) Close() {
	for _, server := range sm.servers {
		server.Close()
	}
}
