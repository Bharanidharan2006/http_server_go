package header

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, []byte("\r\n"))

	if index == -1 {
		return 0, false, nil
	}

	read := 0

	headerLines := bytes.Split(data, []byte("\r\n"))
	
	for i := 0; i < len(headerLines) - 1 ; i++ {
		element := headerLines[i]
		length := len(element)
		line := string(element)
		if line == "" {
			return read, done, nil 
		}

		parts := strings.SplitN(line, ":", 2)

		parts[0] = strings.TrimLeft(parts[0], " ")
		parts[1] = strings.Trim(parts[1], " ")

		if strings.Contains(parts[0], " ") {
			return 0, false, errors.New("Header is not in the correct format")
		}

		read += (length + len([]byte("\r\n")))

		h[parts[0]] = parts[1]
		
	}

	return read, false, nil
}