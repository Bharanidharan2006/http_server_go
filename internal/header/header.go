package header

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Get(key string) string {
	return h.Get(key)
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	
	read := 0
	
	for {
		index := bytes.Index(data, []byte("\r\n"))
		if index == -1 {
			return read, false, nil
		}

		if index == 0 {
			return read + len([]byte("\r\n")), true, nil
		}
		
		line := string(data[:index])
		read += index + len([]byte("\r\n"))

		parts := strings.SplitN(line, ":", 2)

		parts[0] = strings.TrimLeft(parts[0], " ")
		parts[1] = strings.Trim(parts[1], " ")

		if strings.Contains(parts[0], " ") {
			return 0, false, errors.New("Header is not in the correct format")
		}

		matched, _ := regexp.MatchString("^[a-zA-Z0-9!#$%&+-.*'^_`~|]+$", parts[0])
		
		if !matched {
			return 0, false, errors.New("Header is not in the correct Format. Use of invalid characters in the header")
		}
		
		parts[0] = strings.ToLower(parts[0])

		val , ok := (*h)[parts[0]]

		if ok {
			newVal := val + ", " + parts[1]
			(*h)[parts[0]] = newVal
		} else {
			(*h)[parts[0]] = parts[1]
		}

		offset := index + len([]byte("\r\n"))

		copy(data, data[offset:])
		
	}
}