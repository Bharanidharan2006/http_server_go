package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	isClosed atomic.Bool
	handler Handler
}

type HandlerError struct {
	Message string
	StatusCode response.StatusCode
}

type Handler func(w *response.Writer, req *request.Request) 

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		return
	}

	s.handler(&response.Writer{Writer:conn }, req)

}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			if s.isClosed.Load() == true {
				return
			} else {
				continue
			}
		}

		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.isClosed.Store(true)
	return nil
}

func Serve(port int, handler Handler) (*Server, error){
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := &Server{Listener: listener, handler: handler} 

	go server.listen()

	return server, nil
}

