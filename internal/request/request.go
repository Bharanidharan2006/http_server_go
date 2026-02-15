package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	isDone uint8
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() (*Request) {
	return  &Request{isDone: 0}
}

func (r *Request) parse(data []byte) (int, error){
	rl, np, err := parseRequestLine(data)

	if err != nil {
		return 0, err
	}

	if np == 0 {
		return 0, nil
	}

	r.RequestLine = *rl
	r.isDone = 1

	return np, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	// Request data can exceed the buffer length 
	buf := make([]byte, 1024)
	readIndex := 0
	for request.isDone == 0 {
		n, err := reader.Read(buf[readIndex:])

		if err != nil {
			return nil, err
		}

		readIndex += n

		parsed, err := request.parse(buf)
		if err != nil {
			return nil, err
		}

		if parsed != 0 {
			copy(buf, buf[parsed: readIndex])
			readIndex -= parsed
		}


	} 
	return request, nil


}

func parseRequestLine(requestLine []byte) (*RequestLine, int, error){

	read := bytes.Index(requestLine, []byte("\r\n"))

	if read == -1 {
		return nil, 0 , nil
	}

	requestLine = requestLine[:read]

	parts := bytes.Split(requestLine, []byte(" "))

	if len(parts) != 3 {
		return  nil,0, errors.New("Request line is not in the correct format")
	}

	method := string(parts[0])
	requestTarget := string(parts[1])
	httpVersion := string(parts[2])

	httpVersionNo := strings.Split(httpVersion, "/")[1]

	if httpVersion != "HTTP/1.1" {
		return nil, 0, errors.New("Http version mismatch")
	}

	return &RequestLine{
		HttpVersion: httpVersionNo,
		RequestTarget: requestTarget,
		Method: method,
	}, read + len([]byte("\r\n")) , nil

}