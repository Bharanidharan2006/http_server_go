package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("error", err.Error())
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("error", err.Error())
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	out := os.Stdout

	for {
		out.Write([]byte(">"))
		switch line,_,  err := reader.ReadLine(); err {
		case nil:
				conn.Write(line)
		case io.EOF:
				return
		default:
				log.Fatal("error", err.Error())
		}
	}

}