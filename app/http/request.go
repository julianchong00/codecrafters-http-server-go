package http

import (
	"bufio"
	"bytes"
	"fmt"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

func ParseRequest(request []byte) (Request, error) {
	reader := bytes.NewReader(request)
	scanner := bufio.NewScanner(reader)

	var req Request
	req.Headers = make(map[string]string)
	firstLine := true
	for scanner.Scan() {
		line := scanner.Bytes()
		if firstLine {
			words := bytes.Split(line, []byte(" "))
			if len(words) != 3 {
				return req, fmt.Errorf("expected start line to be 3 tokens, but was %d", len(words))
			}
			req.Method = string(words[0])
			req.Path = string(words[1])
			req.Version = string(words[2])
			firstLine = false
			continue
		}
	}

	return req, nil
}
