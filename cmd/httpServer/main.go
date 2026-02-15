package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{
				Message: "Not my problem. You have to solve your own problems\n",
				StatusCode: response.BadRequest,
			}
		}
		if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{
				Message: "Sorry an error occured on my side\n",
				StatusCode: response.InternalServerError,
			}
		}
		
		w.Write([]byte("All good, frfr\n"))
		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
