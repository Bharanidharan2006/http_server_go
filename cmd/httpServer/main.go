package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"httpfromtcp/internal/header"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		if req.RequestLine.RequestTarget == "/" {
			headers := header.NewHeaders()

			str := "<html> <head> <title>200 OK</title> </head> <body> <h1>Success!</h1> <p>Your request was an absolute banger.</p> </body></html>"

			headers["Content-Type"] = "text/html"
			headers["Connection"] = "close"
			headers["Content-Length"] = strconv.Itoa(len([]byte(str)))

			w.WriteStatusLine(response.Ok)
			w.WriteHeaders(headers)
			w.WriteBody([]byte(str))
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			headers := header.NewHeaders()
			headers["Content-Type"] = "text/plain"
			headers["Connection"] = "close"
			headers["Transfer-Encoding"] = "chunked"
			headers["Trailer"] = "X-Content-SHA256, X-Content-Length"
			x := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

			res, _ := http.Get(fmt.Sprintf("https://httpbin.org/%s", x))

			buff := make([]byte, 1024)

			fullBuff := make([]byte, 1024)
			fullLength := 0

			for {
				n, err := res.Body.Read(buff)

				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
				}

				fullLength += n
				fullBuff = append(fullBuff, buff[:n]...)

				w.WriteChunckedBody(buff[:n], n)
			}

			trailingHeaders := header.NewHeaders()
			trailingHeaders["X-Content-Length"] = strconv.Itoa(fullLength)
			trailingHeaders["X-Content-SHA256"] = fmt.Sprintf("%x", sha256.Sum256(fullBuff))
			w.WriteTrailers(trailingHeaders)

		} else if req.RequestLine.RequestTarget == "/video" {
			headers := header.NewHeaders()
			headers["Content-Type"] = "video/mp4"
			headers["Connection"] = "close"
			video, err := os.ReadFile("assets/vim.mp4")
			
			if err != nil {
				return
			}
			headers["Content-Length"] = fmt.Sprintf("%d", len(video))

			w.WriteStatusLine(response.Ok)
			w.WriteHeaders(headers)
			w.WriteBody(video)
		}
	
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
