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
	writer io.Writer
}

const (
	Ok StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	err := WriteStatusLine(w.writer, statusCode)

	if err != nil {
		return err
	}

	return nil
} 

func (w *Writer) WriteHeaders(headers header.Headers) error {
	err := WriteHeaders(w.writer, headers)

	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) writeBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	if err != nil {
		return 0, err
	}

	return n, nil
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