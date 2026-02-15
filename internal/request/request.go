package request

import (
	"bytes"
	"errors"
	"httpfromtcp/internal/header"
	"io"
	"strconv"
	"strings"
)

type parserState string

const (
    InitialState parserState = "initial"
    RequestLineParsedState parserState = "request-line-parsed"
    HeaderParsedState parserState = "header-parsed"
    DoneState parserState = "done"
)

type Request struct {
    RequestLine RequestLine
    Headers header.Headers
    Body []byte
    state parserState
}

type RequestLine struct {
    HttpVersion   string
    RequestTarget string
    Method        string
}

func newRequest() (*Request) {
    return  &Request{state: InitialState, Headers: header.NewHeaders(), Body: make([]byte, 0)}
}

func (r *Request) parse(data []byte) (int, error){

    switch r.state {
    case InitialState:
        rl, np, err := parseRequestLine(data)

        if err != nil {
            return 0, err
        }

		if np == 0 {
			return 0, nil
		}

        r.RequestLine = *rl

        r.state = RequestLineParsedState

		return np, nil

    case RequestLineParsedState:
        np, done, err := r.Headers.Parse(data)

        if err != nil {
    	      return 0, err
        }

        if done {
            r.state = HeaderParsedState
        }

		return np, nil

    	default:
       	 return 0, nil
		}

	
}

func (r * Request) done() bool {
    return r.state == DoneState
}

func (r * Request) headerParsed() bool {
    return r.state == HeaderParsedState
}
func RequestFromReader(reader io.Reader) (*Request, error) {
    request := newRequest()
    // Request data can exceed the buffer length
    buf := make([]byte, 1024)
    readIndex := 0
    for !request.headerParsed() {
        n, err := reader.Read(buf[readIndex:])

        if err != nil {
            if errors.Is(err, io.EOF) {
				return nil, errors.New("connection closed before headers were fully received")
			}
			return nil, err 
		}

        readIndex += n

		for {
			parsed, err := request.parse(buf[:readIndex])

        	if err != nil {
           	 	return nil, err
        	}

			if parsed == 0 {
				break
			}

        	
            copy(buf, buf[parsed: readIndex])
            readIndex -= parsed

			if request.headerParsed() {
				break
			}

		}

        
    }

    if _, ok := request.Headers["content-length"] ; !ok {
		request.state = DoneState
        return request, nil
    }

    cl, err := strconv.Atoi(request.Headers["content-length"])

    if err != nil {
        return nil, err
    }

	request.Body = append(request.Body, buf[:readIndex]...)
	bodyBuf := make([]byte, 1024)
	read := len(request.Body)
	if len(request.Body) >= cl {
		if len(request.Body) > cl {
            return nil, errors.New("Body length is more than that specified in the Content Length")
        }
		request.state = DoneState
	}
    for !request.done() {
        n, err := reader.Read(bodyBuf)
        if err != nil {
            if errors.Is(err, io.EOF){
				if read < cl {
					return nil, errors.New("Body length is less than the Content Length")
				}
                break
            }

            return nil, err
        }

		read += n

        

        request.Body = append(request.Body, bodyBuf[:n]...)

		if len(request.Body) >= cl {
			if len(request.Body) > cl {
            return nil, errors.New("Body length is more than that specified in the Content Length")
        	}
			request.state = DoneState
			break
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