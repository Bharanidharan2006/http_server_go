package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/header"
	"io"
	"strconv"
)

type StatusCode int

type Writer struct {
	Writer io.Writer
}

const (
	Ok StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	err := WriteStatusLine(w.Writer, statusCode)

	if err != nil {
		return err
	}

	return nil
} 

func (w *Writer) WriteHeaders(headers header.Headers) error {
	err := WriteHeaders(w.Writer, headers)

	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Writer.Write(p)

	if err != nil {
		return 0, err
	}

	return n, nil
}


func (w *Writer) WriteChunckedBody(p []byte, n int) (int, error) {
	read := 0

	for i := 0; i <= (n/16); i++ {
		var buff []byte

		if i == (n/16){
			buff = p[i*16:]
		} else {
			buff = p[i*16:(i*16)+16]
		}

		str := fmt.Sprintf("%x\r\n", len(buff))

		n1, _ := w.Writer.Write([]byte(str))
		n2, err := w.Writer.Write(buff)

		read += (n1 + n2)

		if err != nil {
			return 0, err
		}
	}

	return read, nil
}

func (w *Writer) WriteChunckedBodyEnd() (int, error) {

	str := fmt.Sprintf("%x\r\n\r\n", 0)

	n, err := w.Writer.Write([]byte(str))

	if err != nil {
		return 0, err
	}

	return n , nil
}

func (w *Writer) WriteTrailers(h header.Headers) error {
	str := fmt.Sprintf("%x\r\n", 0)
	w.Writer.Write([]byte(str))

	for k, v := range h {
		w.Writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}

	w.Writer.Write([]byte("\r\n"))

	return nil
}


func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	str := ""

	switch statusCode {
		case Ok:
			str += "HTTP/1.1 200 OK\r\n"
		case BadRequest:
			str += "HTTP/1.1 400 Bad Request\r\n"
		case InternalServerError:
			str += "HTTP/1.1 500 Internal Server Error\r\n"
		default:
			return errors.New("Invalid Status Code")
	}

	_, err := w.Write([]byte(str))

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) header.Headers {
	defaultHeaders := header.NewHeaders()
	defaultHeaders["Content-Type"] = "text/plain"
	defaultHeaders["Connection"] = "close"
	defaultHeaders["Content-Length"] = strconv.Itoa(contentLen)
	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers header.Headers) error {
	for k, v := range headers {
		str := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(str))

		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}