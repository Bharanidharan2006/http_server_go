package main

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/request"
	"io"
	"net"
	"strings"
)

func main(){
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}

		// linesChan := getLinesChannel(conn)

		// for line := range linesChan {
		// 	fmt.Println(line)
		// }

		// conn.Close()

		req, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("Request Line:\n\t-Method: %s\n\t-Request Target: %s\n\t-Http Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
	}
}


func getLinesChannel(f io.ReadCloser) <-chan string {
	
	lineChan := make(chan string)
	go func() {
		defer f.Close()
		defer close(lineChan)
		var currline string 
		for {
		b := make([]byte, 8, 8)
		n, err := f.Read(b)

		if err != nil {
			if currline != "" {
				lineChan <- currline
			} 
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("%s\n", err.Error())
			return
		}

		currline += string(b[:n])

		if strings.Contains(currline, "\n") {
			parts := strings.Split(currline, "\n")
			lineChan <- parts[0]

			if len(parts) > 1 {
				currline = parts[1]
			} else {
				currline = ""
			}

		}
	}	
	
	}()

	return lineChan
}