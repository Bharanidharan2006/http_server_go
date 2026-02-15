package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (s *Server) handle(conn net.Conn) {
	req, err := request.RequestFromReader(conn)
	
	if err != nil {
		return
	}

	buffer := bytes.Buffer{}

	handlerError := s.handler(&buffer, req)

	if handlerError != nil {
		b := []byte(handlerError.Message)
		headers := response.GetDefaultHeaders(len(b))
		response.WriteStatusLine(conn, handlerError.StatusCode)
		response.WriteHeaders(conn, headers)
		conn.Write(b)
		conn.Close()
	} else {
		headers := response.GetDefaultHeaders(buffer.Len())
		response.WriteStatusLine(conn, response.Ok)
		response.WriteHeaders(conn, headers)
		conn.Write(buffer.Bytes())
		conn.Close()
	}

	

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

