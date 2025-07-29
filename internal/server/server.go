package server

import (
	"bufio"
	"fmt"
	"gofkgo/internal/ports"
	"io"
	"log"
	"net"
)

type Server struct {
	Addr    string
	Handler ports.ProtocolPort
}

func NewServer(addr string, h ports.ProtocolPort) ports.Server {
	return &Server{Addr: addr, Handler: h}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("Server is listening on port %s...", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		msg, err := s.Handler.ReadFullMessage(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
			}
			return
		}
		if err := s.Handler.ProtocolHandler(msg, c); err != nil {
			log.Printf("handler error: %v", err)
			return
		}
	}
}
