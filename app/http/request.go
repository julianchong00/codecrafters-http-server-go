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
	var req Request
	req.Headers = make(map[string]string)

	parts := bytes.Split(request, []byte(CRLF+CRLF))
	if len(parts) < 2 {
		return req, fmt.Errorf("expected 2 parts in request, but was %d", len(parts))
	}

	// Read request method, path, version, and headers
	reader := bytes.NewReader(parts[0])
	scanner := bufio.NewScanner(reader)
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

		words := bytes.Split(line, []byte(": "))
		if len(words) == 1 {
			continue
		} else if len(words) != 2 {
			return req, fmt.Errorf("expected header line to be 2 tokens, but was %d", len(words))
		}
		req.Headers[string(words[0])] = string(words[1])
	}

	// Read request body
	reader = bytes.NewReader(parts[1])
	scanner = bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		req.Body = append(req.Body, line...)
	}

	return req, nil
}
